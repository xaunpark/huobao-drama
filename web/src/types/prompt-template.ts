export interface PromptTemplate {
  id: number
  name: string
  description?: string
  prompts: PromptTemplatePrompts
  created_at: string
  updated_at: string
}

export interface PromptTemplatePrompts {
  storyboard_breakdown?: string
  character_extraction?: string
  scene_extraction?: string
  prop_extraction?: string
  script_outline?: string
  script_episode?: string
  image_first_frame?: string
  image_key_frame?: string
  image_last_frame?: string
  image_action_sequence?: string
  video_constraint?: string
  style_prompt?: string
  visual_unit_breakdown?: string  // Voice-over AI Director shot planning
  narrative_music_dna?: string    // Music-specific style DNA for narrative_mv mode
}

export interface CreatePromptTemplateRequest {
  name: string
  description?: string
  prompts: PromptTemplatePrompts
}

export interface UpdatePromptTemplateRequest {
  name?: string
  description?: string
  prompts?: PromptTemplatePrompts
}

// Prompt types grouped by category for Tab UI
export const PROMPT_TYPE_GROUPS = [
  {
    label: '📝 Kịch bản',
    key: 'script',
    types: [
      { key: 'script_outline', label: 'Tạo dàn ý (Outline)' },
      { key: 'script_episode', label: 'Tạo kịch bản chia tập' },
    ]
  },
  {
    label: '🎭 Trích xuất',
    key: 'extraction',
    types: [
      { key: 'character_extraction', label: 'Trích xuất nhân vật' },
      { key: 'scene_extraction', label: 'Trích xuất bối cảnh' },
      { key: 'prop_extraction', label: 'Trích xuất đạo cụ' },
    ]
  },
  {
    label: '🎬 Phân cảnh',
    key: 'storyboard',
    types: [
      { key: 'storyboard_breakdown', label: 'Phân rã Storyboard' },
      { key: 'visual_unit_breakdown', label: 'AI Director (Voice-over)' },
    ]
  },
  {
    label: '🖼️ Hình ảnh',
    key: 'image',
    types: [
      { key: 'image_first_frame', label: 'Prompt ảnh: Đầu cảnh' },
      { key: 'image_key_frame', label: 'Prompt ảnh: Cảnh đinh' },
      { key: 'image_last_frame', label: 'Prompt ảnh: Cuối cảnh' },
      { key: 'image_action_sequence', label: 'Prompt ảnh: Dải hành động 1×3' },
    ]
  },
  {
    label: '🎥 Video',
    key: 'video',
    types: [
      { key: 'video_constraint', label: 'Ràng buộc Video' },
    ]
  },
  {
    label: '🎨 Phong cách',
    key: 'style',
    types: [
      { key: 'style_prompt', label: 'Prompt phong cách chung' },
      { key: 'narrative_music_dna', label: '🎵 Music DNA (Narrative MV)' },
    ]
  }
] as const
