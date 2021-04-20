package app

import (
	contentserviceapi "contentservice/api/contentservice"
	userserviceapi "contentservice/api/userservice"
)

func RunTests(contentServiceClient contentserviceapi.ContentServiceClient, userServiceClient userserviceapi.UserServiceClient) {
	ContentTests(contentServiceClient, userServiceClient)
}
