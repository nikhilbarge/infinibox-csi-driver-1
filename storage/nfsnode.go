package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (nfs *nfsstorage) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, time.Now().String()+"---  NodeStageVolume not implemented")
}
func (nfs *nfsstorage) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, time.Now().String()+"---  NodeUnstageVolume not implemented")
}

func (nfs *nfsstorage) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Debug("NodePublishVolume")
	log.Infof("NodePublishVolume volumeID %v", req.GetVolumeId())
	log.Infof("NodePublishVolume GetVolumeContext====>>> %v", req.GetVolumeContext())
	log.Infof("NodePublishVolume  GetPublishContext====>>> %v", req.GetPublishContext())
	log.Infof("Node IP address %v", nfs.cs.nodeIPAddress)
	targetPath := req.GetTargetPath()
	notMnt, err := nfs.mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(targetPath, 0750); err != nil {
				log.Errorf("Error while mkdir %v", err)
				return nil, err
			}
			notMnt = true
		} else {
			log.Errorf("IsLikelyNotMountPint method error  %v", err)
			return nil, err
		}
	}
	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}
	mountOptions := []string{}
	configMountOptions := req.GetVolumeContext()["nfs_mount_options"]
	if configMountOptions == "" {
		configMountOptions = MountOptions
	}
	for _, option := range strings.Split(configMountOptions, ",") {
		if option != "" {
			mountOptions = append(mountOptions, option)
		}
	}
	if req.GetReadonly() {
		mountOptions = append(mountOptions, "ro")
	}
	log.Debug("Mount options  %v", mountOptions)
	sourceIP := req.GetVolumeContext()["ipAddress"]
	ep := req.GetVolumeContext()["volPathd"]
	source := fmt.Sprintf("%s:%s", sourceIP, ep)
	log.Debug("Mount sourcePath %v, tagetPath %v", source, targetPath)
	err = nfs.mounter.Mount(source, targetPath, "nfs", mountOptions)
	if err != nil {
		log.Errorf("fail to mount source path '%s' : %s", source, err)
		return nil, status.Errorf(codes.Internal, "Failed to mount target path '%s': %s", targetPath, err)
	}
	return &csi.NodePublishVolumeResponse{}, nil
}

func (nfs *nfsstorage) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	log.Debug("NodeUnpublishVolume Method")
	log.Infof("** NodeUnpublishVolume Node IP address %v", nfs.cs.nodeIPAddress)
	targetPath := req.GetTargetPath()
	notMnt, err := nfs.mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warnf("mount point '%s' already doesn't exist: '%s', return OK", targetPath, err)
			return &csi.NodeUnpublishVolumeResponse{}, nil
		}
		return nil, err
	}
	if notMnt {
		if err := os.Remove(targetPath); err != nil {
			log.Errorf("Remove target path error: %s", err.Error())
		}
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}
	if err := nfs.mounter.Unmount(targetPath); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to unmount target path '%s': %s", targetPath, err)
	}
	if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
		return nil, status.Errorf(codes.Internal, "Cannot remove unmounted target path '%s': %s", targetPath, err)
	}

	err = DeleteExportRule(req.GetVolumeId(), nfs.cs.nodeIPAddress)
	if err != nil {
		log.Errorf("fail to delete export rule %v", err)
		return nil, status.Errorf(codes.Internal, "fail to delete export rule  %s", err)
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (nfs *nfsstorage) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{}, nil
}

func (nfs *nfsstorage) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return nil, nil

}

func (nfs *nfsstorage) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, time.Now().String()+"---  NodePublishVolume not implemented")

}

func (nfs *nfsstorage) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, time.Now().String()+"---  NodePublishVolume not implemented")
}
