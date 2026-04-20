import type {
  GenerateVideoRequest,
  VideoGeneration,
  VideoGenerationListParams
} from '../types/video'
import request from '../utils/request'

export const videoAPI = {
  generateVideo(data: GenerateVideoRequest) {
    return request.post<VideoGeneration>('/videos', data)
  },

  generateFromImage(imageGenId: number) {
    return request.post<VideoGeneration>(`/videos/image/${imageGenId}`)
  },

  batchGenerateForEpisode(episodeId: number) {
    return request.post<VideoGeneration[]>(`/videos/episode/${episodeId}/batch`)
  },

  getVideoGeneration(id: number) {
    return request.get<VideoGeneration>(`/videos/${id}`)
  },
  
  getVideo(id: number) {
    return request.get<VideoGeneration>(`/videos/${id}`)
  },

  listVideos(params: VideoGenerationListParams) {
    return request.get<{
      items: VideoGeneration[]
      pagination: {
        page: number
        page_size: number
        total: number
        total_pages: number
      }
    }>('/videos', { params })
  },

  deleteVideo(id: number) {
    return request.delete(`/videos/${id}`)
  },

  upscaleVideo(videoGenId: number) {
    return request.post<{ status: string; message: string }>(`/videos/${videoGenId}/upscale`)
  },
  
  resetVideoStatus(videoGenId: number) {
    return request.post<{ status: string; message: string }>(`/videos/${videoGenId}/reset-status`)
  },

  reviewVideo(videoGenId: number) {
    return request.post<{ task_id: string; status: string; message: string }>(`/videos/${videoGenId}/review`)
  },

  getVideoReview(videoGenId: number) {
    return request.get<any>(`/videos/${videoGenId}/review`)
  }
}
