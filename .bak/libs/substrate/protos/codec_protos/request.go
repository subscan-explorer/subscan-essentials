package codec_protos

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"subscan-end/utiles"
	"time"
)

func DecodeExtrinsic(extrinsics string, version int) (string, error) {
	conn, err := utiles.GrpcFromPythonClient()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := NewToolsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.DecodeExtrinsic(ctx, &ExtrinsicRequest{Message: extrinsics, MetadataVersion: int32(version)})
	if err != nil {
		return "", err
	}
	return r.Message, nil
}

func DecodeEvent(event string, version int) (string, error) {
	conn, err := utiles.GrpcFromPythonClient()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := NewToolsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.DecodeEvent(ctx, &EventRequest{Message: event, MetadataVersion: int32(version)})
	if err != nil {
		return "", err
	}
	return r.Message, nil
}

type DecoderLog struct {
	Index string                 `json:"index"`
	Type  string                 `json:"type"`
	Value map[string]interface{} `json:"value"`
}

func DecodeLog(log string, version int) (string, error) {
	conn, err := utiles.GrpcFromPythonClient()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := NewToolsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.DecodeLog(ctx, &LogRequest{Message: log, MetadataVersion: int32(version)})
	if err != nil {
		return "", err
	}
	return r.Message, nil
}

func DecodeStorage(storage, decodeType string) (string, error) {
	if regexp.MustCompile("^[0-9a-fA-F]+$").MatchString(utiles.TrimHex(storage)) == false {
		return "", errors.New(fmt.Sprintf("%s not hex string", storage))
	}
	conn, err := utiles.GrpcFromPythonClient()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := NewToolsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.DecodeStorage(ctx, &StorageRequest{Message: storage, DecoderType: decodeType})
	if err != nil {
		return "", err
	}
	return r.Message, nil
}
