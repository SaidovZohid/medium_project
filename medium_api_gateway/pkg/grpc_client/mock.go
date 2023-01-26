package grpc_client

import (
	pbp "gitlab.com/medium-project/medium_api_gateway/genproto/post_service"
	pbu "gitlab.com/medium-project/medium_api_gateway/genproto/user_service"
)

func (g *GrpcClient) SetUserService(u pbu.UserServiceClient) {
	g.connections["user_service"] = u
}

func (g *GrpcClient) SetAuthService(u pbu.AuthServiceClient) {
	g.connections["auth_service"] = u
}

func (g *GrpcClient) SetPostService(p pbp.PostServiceClient) {
	g.connections["post_service"] = p
}

func (g *GrpcClient) SetCategoryService(c pbp.CategoryServiceClient) {
	g.connections["category_service"] = c
}

func (g *GrpcClient) SetCommentService(c pbp.CommentServiceClient) {
	g.connections["comment_service"] = c
}

func (g *GrpcClient) SetLikeService(l pbp.LikeServiceClient) {
	g.connections["like_service"] = l
}
