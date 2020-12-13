package storage

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultBucketSuffix  = "storage"
	defaultTempDirPrefix = "storage"

	envPrefix           = "storage"
	envBucketKey        = "bucket"
	envProjectKey       = "project_id"
	envTempDirPrefixKey = "temp_dir_prefix"

	headerUserInfo = "X-Endpoint-API-UserInfo"
)

type Options struct {
	Bucket    string
	ProjectId string
	TempDir   string
}

type AuthUserInfo struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Hd            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
}

type option func(*Options)

func withBucket(bucket string) option {
	return func(o *Options) {
		o.Bucket = bucket
	}
}

func withProject(project string) option {
	return func(o *Options) {
		o.ProjectId = project
	}
}

func withTempDir(tmpdir string) option {
	return func(o *Options) {
		o.TempDir = tmpdir
	}
}

func Bootstrap() (*Service, error) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	optf := make([]option, 0)

	prj, err := reqEnvString(envProjectKey)
	if err != nil {
		return nil, err
	}
	optf = append(optf, withProject(prj))

	bkt := viper.GetString(envBucketKey)
	if bkt == "" {
		bkt = fmt.Sprintf("%s-%s", prj, defaultBucketSuffix)
	}
	optf = append(optf, withBucket(bkt))

	// Create temporary directory
	suffix := viper.GetString(envTempDirPrefixKey)
	if suffix == "" {
		suffix = defaultTempDirPrefix
	}
	pattern := fmt.Sprintf("*-%s", suffix)
	dir, err := ioutil.TempDir("", pattern)
	if err != nil {
		return nil, err
	}
	optf = append(optf, withTempDir(dir))

	opts := &Options{}
	for _, f := range optf {
		f(opts)
	}

	svc, err := NewStorageService(opts)
	if err != nil {
		return nil, err
	}

	if err := svc.validate(bkt); err != nil {
		return nil, err
	}

	return svc, nil
}

func reqEnvString(key string) (string, error) {
	val := viper.GetString(key)
	if val == "" {
		return val, fmt.Errorf("env=[%s] required", strings.ToUpper(key))
	}
	return val, nil
}

func hashFileName(s string) (b []byte) {
	h := sha1.New()
	h.Write([]byte(s))
	b = h.Sum(nil)
	return
}

func ExtractUserInfo(c *gin.Context) (*AuthUserInfo, error) {
	str := c.Request.Header.Get(headerUserInfo)
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	uinfo := &AuthUserInfo{}
	err = json.Unmarshal(data, uinfo)
	if err != nil {
		return nil, err
	}

	return uinfo, nil
}
