package storage

import (
	"fmt"
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
)

type Options struct {
	Bucket    string
	ProjectId string
	TempDir   string
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

	// if err := svc.validate(bkt); err != nil {
	// 	return nil, err
	// }

	return svc, nil
}

func reqEnvString(key string) (string, error) {
	val := viper.GetString(key)
	if val == "" {
		return val, fmt.Errorf("env=[%s] required", strings.ToUpper(key))
	}
	return val, nil
}
