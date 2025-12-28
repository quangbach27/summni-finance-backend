package assetsource

import "context"

type Repository interface {
	GetByID(ctx context.Context, id ID) (*AssetSource, error)
	Create(ctx context.Context, assetSourceList []*AssetSource) error
}
