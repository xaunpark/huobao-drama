import type {
  PromptTemplate,
  PromptTemplatePrompts,
  CreatePromptTemplateRequest,
  UpdatePromptTemplateRequest
} from '../types/prompt-template'
import request from '../utils/request'

export const promptTemplateAPI = {
  list() {
    return request.get<PromptTemplate[]>('/prompt-templates')
  },

  get(id: number) {
    return request.get<PromptTemplate>(`/prompt-templates/${id}`)
  },

  create(data: CreatePromptTemplateRequest) {
    return request.post<PromptTemplate>('/prompt-templates', data)
  },

  update(id: number, data: UpdatePromptTemplateRequest) {
    return request.put<PromptTemplate>(`/prompt-templates/${id}`, data)
  },

  delete(id: number) {
    return request.delete(`/prompt-templates/${id}`)
  },

  duplicate(id: number) {
    return request.post<PromptTemplate>(`/prompt-templates/${id}/duplicate`)
  },

  getDefaults() {
    return request.get<PromptTemplatePrompts>('/prompt-templates/defaults')
  }
}
