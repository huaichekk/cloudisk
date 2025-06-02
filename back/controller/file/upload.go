package file

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"io"
	"netdisk/commen/database"
	"netdisk/commen/redis"
	"netdisk/models"
	"netdisk/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	UPLOAD_TEMP_PATH = "./upload/temp/"
	UPLOAD_PATH      = "./upload/"
)

type FileMetaData struct {
	Name     string `json:"file_name"`
	Size     int    `json:"file_size"`
	Hash     string `json:"file_hash"`
	FatherID int    `json:"father_id"`
}

const CHUNK_SIZE = 1 * 1024 * 1024 //1M

// 上传初始化，若文件hash存在触发秒传，若redis中有上传进度信息，则断点续传，否则只进行初始化
func UploadInitHandle() gin.HandlerFunc {
	//Post http:127.0.0.1:8080/api/file/init
	return func(c *gin.Context) {
		var metadata FileMetaData
		if err := c.BindJSON(&metadata); err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		//检查本地有没有已经存储了该文件
		var file models.File
		if err := database.DB().Where("hash = ?", metadata.Hash).First(&file).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { //服务器本地没有存储了该文件
				//查看redis中有没有相关信息，如果有触发断点续传
				_, err := redis.GetRedis().HGet(context.Background(), metadata.Hash, "chunk_num").Result()
				if errors.Is(err, redis2.Nil) { //没找到上传进度记录
					redis.GetRedis().HSet(context.Background(), metadata.Hash, []string{
						"chunk_num", strconv.Itoa(getChunkNum(CHUNK_SIZE, metadata.Size)),
						"file_name", metadata.Name,
						"father_id", strconv.Itoa(metadata.FatherID),
						"file_size", strconv.Itoa(metadata.Size),
					})
					// 然后设置 TTL
					expireTime := 24 * time.Hour
					err = redis.GetRedis().Expire(context.Background(), metadata.Hash, expireTime).Err()
					if err != nil {
						utils.ServerFailWithMessage(c, err.Error())
						return
					}
					resp := map[string]interface{}{
						"has_upload": false,
						"chunk_size": CHUNK_SIZE,
						"chunk_num":  getChunkNum(CHUNK_SIZE, metadata.Size),
					}
					utils.OkWithData(c, resp)
					return
				} else { //断点续传,将已经上传的分片返回给前端
					//--------------------------------------------------------------------------------------------------
					hasChunk := make([]string, 0)
					if value, err := redis.GetRedis().HGetAll(context.Background(), metadata.Hash).Result(); err != nil {
						utils.ServerFailWithMessage(c, err.Error())
						return
					} else {
						for k, v := range value {
							if strings.HasPrefix(k, "chunk_idx") {
								hasChunk = append(hasChunk, v)
							}
						}
						resp := map[string]interface{}{
							"has_upload": true,
							"chunk_size": CHUNK_SIZE,
							"chunk_num":  getChunkNum(CHUNK_SIZE, metadata.Size),
							"has_chunk":  hasChunk,
						}
						utils.OkWithData(c, resp)
						return
					}
				}
			} else { //其他错误
				utils.ServerFailWithMessage(c, err.Error())
				return
			}
		}

		//--------------------------------------------------------------------------------------------------------------
		//服务器存储了该文件,只需添加用户文件表到全局文件的对应关系即可，触发秒传
		var userFile models.UserFile
		userID, err := utils.GetUserID(c)
		if err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		userFile.UserID = *userID
		userFile.FileID = &file.ID
		userFile.Type = 1
		userFile.FatherID = metadata.FatherID
		if err := database.DB().Create(&userFile).Error; err != nil {
			utils.ServerFailWithMessage(c, err.Error())
			return
		}
		utils.OkWithMessage(c, "触发秒传，上传成功")
	}
}

func getChunkNum(chunkSize, fileSize int) int {
	mod := fileSize % chunkSize
	if mod == 0 {
		return fileSize / chunkSize
	} else {
		return (fileSize / chunkSize) + 1
	}
}

// 上传分片
func UploadChunkHandle() gin.HandlerFunc {
	//Post http:127.0.0.1:8080/api/files/chunk
	return func(c *gin.Context) {
		idx := c.PostForm("chunk_idx")
		file, err := c.FormFile("file")
		if err != nil {
			utils.ServerFailWithMessage(c, err.Error())
			return
		}
		fileHash := c.PostForm("file_hash")
		//创建目录  	./upload/temp/bdhasbjdsv
		if err = os.MkdirAll(UPLOAD_TEMP_PATH+fileHash, 0755); err != nil {
			utils.ServerFailWithMessage(c, "创建分片目录失败:"+err.Error())
			return
		}
		//保存分片   ./upload/temp/bdhasbjdsv/1
		if err = c.SaveUploadedFile(file, UPLOAD_TEMP_PATH+fileHash+"/"+idx); err != nil {
			utils.ServerFailWithMessage(c, "保存分片文件失败:"+err.Error())
			return
		}
		//保存上传进度
		if err = redis.GetRedis().HSet(context.Background(), fileHash, []string{"chunk_idx" + idx, idx}).Err(); err != nil {
			utils.ServerFailWithMessage(c, "Redis保存分片信息失败:"+err.Error())
			return
		}
		utils.OkWithData(c, map[string]string{
			"chunk_idx": idx,
		})
		return
	}
}

// 合并分片
func MergeHandle() gin.HandlerFunc {
	//Post http:127.0.0.1:8080/api/files/merge
	return func(c *gin.Context) {
		var req struct {
			FileHash string `json:"file_hash"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		result, err := redis.GetRedis().HGetAll(context.Background(), req.FileHash).Result()
		if err != nil {
			utils.ServerFailWithMessage(c, err.Error())
			return
		}
		chunkNum, _ := strconv.Atoi(result["chunk_num"])
		fileName := result["file_name"]
		fileSize, _ := strconv.Atoi(result["file_size"])
		fatherID, _ := strconv.Atoi(result["father_id"])
		hasUpload := make([]string, 0)
		for k, v := range result {
			if strings.HasPrefix(k, "chunk_idx") {
				hasUpload = append(hasUpload, v)
			}
		}
		if chunkNum != len(hasUpload) {
			utils.ClientFailWithMessage(c, "分片数量不等于预期数量,上传失败")
			return
		}
		//合并分片
		if err = mergeChunks(req.FileHash, chunkNum); err != nil {
			utils.FailWithData(c, 500, map[string]interface{}{
				"has_upload": hasUpload,
			}, "合并分片失败:"+err.Error())
			return
		}
		//合并成功
		f := &models.File{
			Name: fileName,
			Size: uint64(fileSize),
			Hash: req.FileHash,
		}
		if err = database.DB().Create(f).Error; err != nil {
			utils.ServerFailWithMessage(c, err.Error())
			return
		}
		id, err := utils.GetUserID(c)
		if err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		uf := &models.UserFile{
			UserID:   *id,
			FileID:   &f.ID,
			Type:     1,
			FatherID: fatherID,
			Name:     fileName,
		}
		if err := database.DB().Create(uf).Error; err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		if err := redis.GetRedis().Del(context.Background(), req.FileHash).Err(); err != nil {
			utils.ClientFailWithMessage(c, err.Error())
			return
		}
		if err := removeTempFile(req.FileHash); err != nil {
			utils.ServerFailWithMessage(c, err.Error())
			return
		}
		utils.OkWithMessage(c, "合并成功")
	}
}

func mergeChunks(fileHash string, totalChunks int) error {
	destPath := filepath.Join(UPLOAD_PATH, fileHash)
	tempDir := filepath.Join(UPLOAD_TEMP_PATH, fileHash)

	// 创建目标文件
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %w", err)
	}
	defer func() {
		_ = destFile.Close()
	}()

	buf := make([]byte, 32*1024)

	// 合并分片
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(tempDir, strconv.Itoa(i))

		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("无法打开分片%d: %w", i, err)
		}

		_, err = io.CopyBuffer(destFile, chunkFile, buf)
		_ = chunkFile.Close()

		if err != nil {
			// 删除不完整的目标文件
			_ = os.Remove(destPath)
			return fmt.Errorf("写入分片%d失败: %w", i, err)
		}
	}

	// 确保数据刷到磁盘
	if err := destFile.Sync(); err != nil {
		_ = os.Remove(destPath)
		return fmt.Errorf("同步文件失败: %w", err)
	}

	return nil
}

func removeTempFile(fileHash string) error {
	tempPath := filepath.Join(UPLOAD_TEMP_PATH, fileHash)
	if err := os.RemoveAll(tempPath); err != nil {
		return err
	}
	return nil
}
