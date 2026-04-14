<template>
  <el-dialog
    v-model="visible"
    :title="$t('professionalEditor.batchR2V.title')"
    width="950px"
    :close-on-click-modal="false"
    class="batch-dialog"
  >
    <div class="batch-container">
      <!-- 顶部配置栏 -->
      <div class="batch-config-bar">
        <el-alert
          :title="$t('professionalEditor.batchR2V.title')"
          type="success"
          :description="$t('professionalEditor.batchR2V.instructions')"
          show-icon
          :closable="false"
          style="margin-bottom: 20px"
        />

        <div class="config-row">
          <div class="config-item">
            <span class="label">{{ $t('professionalEditor.batch.videoModel') }}</span>
            <el-select v-model="selectedVideoModel" :placeholder="$t('video.selectVideoModel')" size="small" style="width: 250px">
              <el-option
                v-for="model in videoModels"
                :key="model.id"
                :label="model.name"
                :value="model.id"
              />
            </el-select>
          </div>
          <div class="config-item">
            <el-tooltip :content="$t('professionalEditor.batchR2V.assetPriority')" placement="top">
              <el-tag type="info" size="small">
                <el-icon style="margin-right: 4px"><InfoFilled /></el-icon>
                {{ $t('professionalEditor.batchR2V.assetPriority') }}
              </el-tag>
            </el-tooltip>
          </div>
        </div>

        <div class="action-row" style="margin-top: 15px">
          <el-button type="primary" :loading="isBatching" @click="startFullBatch">
            <el-icon><MagicStick /></el-icon> {{ $t('professionalEditor.batchR2V.runAll') }}
          </el-button>
          <el-button-group>
            <el-button :disabled="isBatching" @click="startStep('prompt')">{{ $t('professionalEditor.batchR2V.onlyVideoPrompt') }}</el-button>
            <el-button :disabled="isBatching" @click="startStep('video')">{{ $t('professionalEditor.batchR2V.onlyVideo') }}</el-button>
          </el-button-group>
          <el-button type="success" plain :disabled="isBatching" @click="startUpscaleAll">
            <el-icon><MagicStick /></el-icon> Upscale All Videos
          </el-button>
          <el-button type="info" plain :loading="isDownloadingZip" @click="downloadAllVideos">
             <el-icon><Download /></el-icon> {{ $t('professionalEditor.batch.downloadAll') }}
          </el-button>
          <el-button type="danger" plain v-if="isBatching" @click="stopBatch">
            {{ $t('professionalEditor.batch.stop') }}
          </el-button>
        </div>
      </div>

      <!-- 分镜列表状态 -->
      <div class="shot-progress-list">
        <el-table :data="storyboards" style="width: 100%" height="450px" v-loading="isBatching && localStoryboards.length === 0">
          <el-table-column :label="$t('professionalEditor.batch.shot')" width="80" property="storyboard_number">
            <template #default="scope">
              #{{ scope.row.storyboard_number }}
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.description')" show-overflow-tooltip>
            <template #default="scope">
              {{ scope.row.action || scope.row.description || $t('professionalEditor.batch.noDescription') }}
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batchR2V.videoPrompt')" width="130">
            <template #default="scope">
              <el-tag :type="getStatusTag(scope.row.id, 'prompt')" size="small">
                {{ getStatusText(scope.row.id, 'prompt') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batchR2V.referenceStatus')" width="120">
            <template #default="scope">
              <el-tag :type="getRefTag(scope.row.id)" size="small">
                {{ getRefText(scope.row.id) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.video')" width="100">
            <template #default="scope">
              <el-tag :type="getStatusTag(scope.row.id, 'video')" size="small">
                {{ getStatusText(scope.row.id, 'video') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('professionalEditor.batch.progress')" width="140">
            <template #default="scope">
              <el-progress 
                :percentage="getProgress(scope.row.id)" 
                :status="getProgressStatus(scope.row.id)"
                :stroke-width="8"
              />
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="visible = false">{{ $t('common.close') || 'Close' }}</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { MagicStick, InfoFilled, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import JSZip from 'jszip'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import { dramaAPI } from '@/api/drama'
import { generateFramePrompt } from '@/api/frame'
import { videoAPI } from '@/api/video'
import { taskAPI } from '@/api/task'
import type { Storyboard } from '@/types/drama'
import { useAISettings } from '@/composables/useAISettings'
import { getImageUrl, getVideoUrl } from '@/utils/image'
import { workerDelay } from '@/utils/worker-timer'

const props = defineProps<{
  modelValue: boolean
  storyboards: Storyboard[]
  episodeId: number
  dramaId: number
  style?: string
  customStyle?: string
  videoModels: any[]
  defaultVideoModel?: string
}>()

const emit = defineEmits(['update:modelValue', 'completed', 'refresh'])
const { t } = useI18n()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const localStoryboards = ref<Storyboard[]>([])
const selectedVideoModel = ref(props.defaultVideoModel || '')
const isBatching = ref(false)
const isDownloadingZip = ref(false)
const shouldStop = ref(false)
const { maxConcurrentThreads } = useAISettings()

// shotId -> { prompt: 'pending'|'loading'|'done'|'error', video: 'pending'|'loading'|'done'|'hd'|'upscaling'|'error', ref: 'ready'|'partial'|'missing', progress: number }
const taskStates = reactive<Record<string, any>>({})

const refreshStoryboardsFromDB = async () => {
  if (!props.episodeId) {
    localStoryboards.value = [...props.storyboards]
    return
  }
  try {
    const res = await dramaAPI.getStoryboards(props.episodeId.toString())
    localStoryboards.value = res?.storyboards || [...props.storyboards]
  } catch (e) {
    localStoryboards.value = [...props.storyboards]
  }
}

const initTaskStates = async () => {
  await refreshStoryboardsFromDB()
  
  for (const sb of localStoryboards.value) {
    // 1. Check prompt
    const hasPrompt = sb.video_prompt && sb.video_prompt_source === 'ai'

    // 2. Check reference assets
    let refStatus = 'missing'
    const sceneImg = getImageUrl((sb as any).background)
    const charImgs = sb.characters?.filter((c: any) => getImageUrl(c)) || []
    const propImgs = (sb as any).props?.filter((p: any) => getImageUrl(p)) || []
    
    const hasScene = !!sceneImg
    const hasAssets = charImgs.length > 0 || propImgs.length > 0

    if (hasScene && hasAssets) {
      refStatus = 'ready'
    } else if (hasScene || hasAssets) {
      refStatus = 'partial'
    }

    // 3. Check video
    let videoState = 'pending'
    let activeVideoId = null
    try {
      const vidRes = await videoAPI.listVideos({ 
        storyboard_id: String(sb.id), 
        page_size: 10 
      })
      // Match direct_r2v mode OR videos that have multiple references (R2V style) but missing mode flag (legacy from earlier today)
      const r2vVideos = vidRes.items?.filter((v: any) => 
        v.generation_mode === 'direct_r2v' || 
        (!v.generation_mode && v.reference_mode === 'multiple')
      ) || []

      const completedVidHd = r2vVideos.find((v: any) => v.status === 'upscaled' || (v.status === 'completed' && v.is_upscaled))
      const upscalingVid = r2vVideos.find((v: any) => v.status === 'upscaling')
      const completedVid = r2vVideos.find((v: any) => v.status === 'completed' && !v.is_upscaled)
      const failedWithBase = r2vVideos.find((v: any) => v.status === 'upscale_failed' || ((v.status === 'failed' || v.status === 'error') && (v.video_url || v.local_path)))

      if (completedVidHd) {
        videoState = 'hd'
        activeVideoId = completedVidHd.id
      } else if (upscalingVid) {
        videoState = 'upscaling'
        activeVideoId = upscalingVid.id
      } else if (completedVid) {
        videoState = 'done'
        activeVideoId = completedVid.id
      } else if (failedWithBase) {
        videoState = 'done'
        activeVideoId = failedWithBase.id
      }
    } catch {
      // Ignore
    }

    let progress = 0
    if (videoState === 'hd' || videoState === 'done') progress = 100
    else if (videoState === 'upscaling') progress = 90
    else if (hasPrompt) progress = 40
    else if (refStatus !== 'missing') progress = 10

    taskStates[sb.id] = {
      prompt: hasPrompt ? 'done' : 'pending',
      video: videoState,
      videoId: activeVideoId,
      ref: refStatus,
      progress: progress
    }
  }
}

const getStatusTag = (id: string, type: string) => {
  const state = taskStates[id]?.[type]
  if (state === 'loading' || state === 'upscaling') return ''
  if (state === 'done') return 'success'
  if (state === 'hd') return 'warning'
  if (state === 'error') return 'danger'
  return 'info'
}

const getStatusText = (id: string, type: string) => {
  const state = taskStates[id]?.[type]
  if (state === 'upscaling') return 'Upscaling...'
  if (state === 'hd') return 'HD'
  if (state === 'loading') return t('professionalEditor.batch.status.loading')
  if (state === 'done') return t('professionalEditor.batch.status.done')
  if (state === 'error') return t('professionalEditor.batch.status.failed')
  return t('professionalEditor.batch.status.pending')
}

const getRefTag = (id: string) => {
  const state = taskStates[id]?.ref
  if (state === 'ready') return 'success'
  if (state === 'partial') return 'warning'
  return 'info'
}

const getRefText = (id: string) => {
  const state = taskStates[id]?.ref
  return t(`professionalEditor.batchR2V.status.${state}`)
}

const getProgress = (id: string) => taskStates[id]?.progress || 0
const getProgressStatus = (id: string) => {
  const s = taskStates[id]
  if (s?.prompt === 'error' || s?.video === 'error') return 'exception'
  if (s?.video === 'done' || s?.video === 'hd') return 'success'
  return undefined
}

const stopBatch = async () => {
  shouldStop.value = true
  isBatching.value = false
  ElMessage.warning(t('professionalEditor.batch.stopping'))
  
  // Clean up 'upscaling' status in base DB records so they revert to 'Ready' (Done)
  for (const id in taskStates) {
    if (taskStates[id].video === 'upscaling') {
      try {
        const vidID = taskStates[id].videoId
        if (vidID) {
          await videoAPI.resetVideoStatus(vidID)
        } else {
          const vidRes = await videoAPI.listVideos({ storyboard_id: id })
          const upscalingVid = vidRes.items?.find((v: any) => v.status === 'upscaling')
          if (upscalingVid) {
            await videoAPI.resetVideoStatus(upscalingVid.id)
          }
        }
        taskStates[id].video = 'done'
      } catch (e) {
        console.error('Failed to reset upscaling status in R2V batch', e)
        taskStates[id].video = 'done'
      }
    } else if (taskStates[id].video === 'loading') {
       taskStates[id].video = 'pending'
    }
  }
}

const processPrompt = async (sb: Storyboard) => {
  taskStates[sb.id].prompt = 'loading'
  taskStates[sb.id].progress = Math.max(taskStates[sb.id].progress, 20)
  try {
    const { task_id } = await generateFramePrompt(Number(sb.id), { frame_type: 'video' })
    
    let result = null
    while (true) {
      if (shouldStop.value) throw new Error('Stopped')
      const task = await taskAPI.getStatus(task_id)
      if (task.status === 'completed') {
        let res = task.result
        if (typeof res === 'string') res = JSON.parse(res)
        result = res.response
        break
      } else if (task.status === 'failed') {
        throw new Error(task.message || 'Prompt extraction failed')
      }
      await workerDelay(2000)
    }

    const finalPrompt = result.single_frame?.prompt || ""
    if (!finalPrompt) throw new Error("AI returned empty video prompt")

    // Video prompt was already saved to DB by backend FramePromptService for type 'video'
    // But we update local ref to avoid UI lag
    sb.video_prompt = finalPrompt
    
    taskStates[sb.id].prompt = 'done'
    taskStates[sb.id].progress = 40
    return finalPrompt
  } catch (e) {
    taskStates[sb.id].prompt = 'error'
    throw e
  }
}

const extractProviderFromModel = (modelName: string): string => {
  const lowerMod = modelName.toLowerCase()
  if (lowerMod.includes("doubao") || lowerMod.includes("seedance")) return "doubao"
  if (lowerMod.includes("runway")) return "runway"
  if (lowerMod.includes("pika")) return "pika"
  if (lowerMod.includes("minimax") || lowerMod.includes("hailuo")) return "minimax"
  if (lowerMod.includes("sora")) return "openai"
  if (lowerMod.includes("kling")) return "kling"
  return "doubao"
}

const processVideo = async (sb: Storyboard) => {
  if (!selectedVideoModel.value) throw new Error('No video model')
  taskStates[sb.id].video = 'loading'
  try {
    const provider = extractProviderFromModel(selectedVideoModel.value)
    
    // Aggregation logic for R2V references: 1 Scene + 2 Characters (fill with Props if needed)
    const referenceImages: string[] = []
    
    // 1. Scene (Priority 1)
    const sceneImg = getImageUrl((sb as any).background)
    if (sceneImg) referenceImages.push(sceneImg)

    // 2. Characters (Priority 2, up to 2)
    const charImgs = sb.characters
      ?.map((c: any) => getImageUrl(c))
      .filter((url): url is string => !!url) || []
    
    charImgs.slice(0, 2).forEach(img => referenceImages.push(img))

    // 3. Props (Refill if slots < 3)
    if (referenceImages.length < 3) {
      const propImgs = (sb as any).props
        ?.map((p: any) => getImageUrl(p))
        .filter((url): url is string => !!url) || []
      
      const needed = 3 - referenceImages.length
      propImgs.slice(0, needed).forEach(img => referenceImages.push(img))
    }

    const result = await videoAPI.generateVideo({
      drama_id: props.dramaId.toString(),
      storyboard_id: Number(sb.id),
      prompt: sb.video_prompt || sb.action || "Cinematic video sequence",
      duration: 5,
      provider: provider,
      model: selectedVideoModel.value,
      reference_mode: 'multiple',
      reference_image_urls: referenceImages,
      generation_mode: 'direct_r2v'
    })
    
    while (true) {
      if (shouldStop.value) throw new Error('Stopped')
      const videoTask = await videoAPI.getVideo(result.id)
      if (videoTask.status === 'completed') {
        taskStates[sb.id].video = 'done'
        taskStates[sb.id].progress = 100
        return videoTask
      } else if (videoTask.status === 'failed') {
        throw new Error(videoTask.error_msg || 'Generation failed')
      }
      await workerDelay(5000)
    }
  } catch (e) {
    taskStates[sb.id].video = 'error'
    throw e
  }
}

const runConcurrently = async <T>(items: T[], limit: number, worker: (item: T) => Promise<void>) => {
  const executing: Promise<void>[] = []
  for (const item of items) {
    if (shouldStop.value) break
    const p = worker(item).finally(() => {
      executing.splice(executing.indexOf(p), 1)
    })
    executing.push(p)
    if (executing.length >= limit) await Promise.race(executing)
  }
  await Promise.all(executing)
}

const startFullBatch = async () => {
  if (!selectedVideoModel.value) {
    ElMessage.warning(t('professionalEditor.batch.selectVideoModelFirst'))
    return
  }
  isBatching.value = true
  shouldStop.value = false
  await initTaskStates()

  // Phase 1: All Prompts
  await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
    if (taskStates[sb.id]?.prompt === 'done') return
    try {
      await processPrompt(sb)
    } catch (e: any) {
      console.error(`Shot ${sb.storyboard_number} prompt failed:`, e)
    }
  })

  // Phase 2: All Videos
  if (!shouldStop.value) {
    await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
      if (taskStates[sb.id]?.prompt !== 'done' || taskStates[sb.id]?.video === 'done' || taskStates[sb.id]?.video === 'hd') return
      try {
        await processVideo(sb)
      } catch (e: any) {
        console.error(`Shot ${sb.storyboard_number} video failed:`, e)
        ElMessage.error(`${t('professionalEditor.batch.shot')} ${sb.storyboard_number} failed: ${e.message}`)
      }
    })
  }

  isBatching.value = false
  if (!shouldStop.value) {
    ElMessage.success(t('professionalEditor.batch.completed'))
    emit('completed')
  }
}

const startStep = async (step: string) => {
  if (step === 'video' && !selectedVideoModel.value) {
    ElMessage.warning(t('professionalEditor.batch.selectVideoModelFirst'))
    return
  }
  isBatching.value = true
  shouldStop.value = false
  await initTaskStates()

  await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
    try {
      if (step === 'prompt' && taskStates[sb.id]?.prompt !== 'done') await processPrompt(sb)
      else if (step === 'video' && taskStates[sb.id]?.video !== 'done' && taskStates[sb.id]?.video !== 'hd') {
        if (taskStates[sb.id]?.prompt !== 'done') await processPrompt(sb)
        await processVideo(sb)
      }
    } catch (e: any) {
      console.error(`Shot ${sb.storyboard_number} step ${step} failed:`, e)
    }
  })

  isBatching.value = false
  emit('completed')
}

const startUpscaleAll = async () => {
  if (isBatching.value) return
  isBatching.value = true
  shouldStop.value = false
  
  try {
    await initTaskStates()
    await runConcurrently(localStoryboards.value, maxConcurrentThreads.value, async (sb) => {
      // Only process shots that have a base video (status 'done') or are already 'upscaling'
      if (!['done', 'upscaling'].includes(taskStates[sb.id]?.video)) return
      
      try {
        const vidRes = await videoAPI.listVideos({ storyboard_id: String(sb.id), page_size: 10 })
        const target = vidRes.items?.find((v: any) => (v.status === 'completed' && !v.is_upscaled) || v.status === 'upscale_failed' || (v.status === 'failed' && (v.video_url || v.local_path)))
        const currentlyUpscaling = vidRes.items?.find((v: any) => v.status === 'upscaling')

        if (target && !currentlyUpscaling) {
          taskStates[sb.id].video = 'upscaling'
          taskStates[sb.id].videoId = target.id
          await videoAPI.upscaleVideo(target.id)
          
          // Poll for completion
          let isFinished = false
          while (!isFinished && !shouldStop.value) {
            await workerDelay(5000)
            const check = await videoAPI.getVideo(target.id)
            if (check.status === 'upscaled' || (check.status === 'completed' && check.is_upscaled)) {
              isFinished = true
              taskStates[sb.id].video = 'hd'
              taskStates[sb.id].progress = 100
            } else if (check.status === 'failed' || check.status === 'upscale_failed') {
              throw new Error(check.error_msg || 'Upscaling failed')
            }
          }
        } else if (currentlyUpscaling) {
          // Wait for the already ongoing upscale
          taskStates[sb.id].video = 'upscaling'
          taskStates[sb.id].videoId = currentlyUpscaling.id
          let isFinished = false
          while (!isFinished && !shouldStop.value) {
            await workerDelay(5000)
            const check = await videoAPI.getVideo(currentlyUpscaling.id)
            if (check.status === 'upscaled' || (check.status === 'completed' && check.is_upscaled)) {
              isFinished = true
              taskStates[sb.id].video = 'hd'
              taskStates[sb.id].progress = 100
            } else if (check.status === 'failed' || check.status === 'upscale_failed') {
              throw new Error(check.error_msg || 'Upscaling failed')
            }
          }
        }
      } catch (err: any) { 
        console.error(`Shot ${sb.storyboard_number} upscale failed:`, err)
        taskStates[sb.id].video = 'error'
      }
    })
  } finally {
    isBatching.value = false
    emit('completed')
  }
}

const downloadAllVideos = async () => {
  if (isBatching.value) return
  isDownloadingZip.value = true
  
  try {
    // 1. Get ALL videos for this drama (not just completed, to catch those that failed upscale but have base URL)
    const allDramaVideos: any[] = []
    let currentPage = 1
    const pageSize = 100
    
    while (true) {
      const vidRes = await videoAPI.listVideos({ 
        drama_id: props.dramaId.toString(), 
        page: currentPage,
        page_size: pageSize
      })
      
      if (vidRes.items && vidRes.items.length > 0) {
        allDramaVideos.push(...vidRes.items)
      }
      
      if (!vidRes.pagination || vidRes.items.length < pageSize || allDramaVideos.length >= vidRes.pagination.total) {
        break
      }
      currentPage++
    }
    
    const usableVideos = allDramaVideos.filter(v => v.video_url || v.local_path)
    
    if (usableVideos.length === 0) {
      ElMessage.warning(t('professionalEditor.batch.noVideosToDownload'))
      return
    }

    // 2. Map storyboard_id to its best video (prefer HD, then latest usable)
    const bestVideos = new Map<number, any>()
    usableVideos.forEach((vid: any) => {
      const sbId = vid.storyboard_id
      if (!sbId) return
      
      const currentBest = bestVideos.get(sbId)
      if (!currentBest) {
        bestVideos.set(sbId, vid)
      } else {
        // Preference logic:
        // 1. HD (upscaled) vs non-HD
        const isVidHd = vid.status === 'upscaled' || vid.is_upscaled
        const isCurrentBestHd = currentBest.status === 'upscaled' || currentBest.is_upscaled
        if (isVidHd && !isCurrentBestHd) {
          bestVideos.set(sbId, vid)
        } else if (isVidHd === isCurrentBestHd) {
          // If both same HD status, pick latest
          if (new Date(vid.created_at) > new Date(currentBest.created_at)) {
            bestVideos.set(sbId, vid)
          }
        }
      }
    })

    // 3. Ensure we cover ALL storyboards if possible
    const itemsToDownload: { fileName: string; url: string }[] = []
    for (const sb of props.storyboards) {
      const vid = bestVideos.get(Number(sb.id))
      if (vid) {
        const sbNum = String(sb.storyboard_number).padStart(3, '0')
        const fileName = `Shot_${sbNum}.mp4`
        const url = getVideoUrl(vid)
        if (url) {
          itemsToDownload.push({ fileName, url })
        }
      }
    }

    if (itemsToDownload.length === 0) {
      ElMessage.warning(t('professionalEditor.batch.noVideosToDownload'))
      return
    }

    ElMessage.info(t('professionalEditor.batch.downloadingZip'))
    
    const zip = new JSZip()
    const downloadPromises = itemsToDownload.map(item => {
      const fullUrl = item.url.startsWith('http') ? item.url : `${window.location.origin}${item.url}`
      return axios.get(fullUrl, { responseType: 'blob' })
        .then(res => {
          zip.file(item.fileName, res.data)
        })
        .catch(err => {
          console.error(`Failed to download ${item.fileName}:`, err)
        })
    })
    
    await Promise.all(downloadPromises)
    
    // 4. Generate and download ZIP
    const content = await zip.generateAsync({ type: 'blob' })
    const zipUrl = URL.createObjectURL(content)
    const link = document.createElement('a')
    link.href = zipUrl
    link.download = `drama_${props.dramaId}_r2v_videos.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(zipUrl)
    
    ElMessage.success(t('professionalEditor.batch.completed'))
  } catch (err: any) {
    console.error('Download all failed:', err)
    ElMessage.error(t('message.operationFailed'))
  } finally {
    isDownloadingZip.value = false
  }
}

watch(() => props.modelValue, (newVal) => {
  if (newVal) initTaskStates()
})

const storyboards = computed(() => {
  return localStoryboards.value.length > 0 ? localStoryboards.value : props.storyboards
})
</script>

<style scoped>
.batch-dialog :deep(.el-dialog__body) {
  padding-top: 10px;
}
.batch-container {
  padding: 0 10px;
}
.batch-config-bar {
  background: var(--el-fill-color-light);
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 24px;
}
.config-row {
  display: flex;
  gap: 30px;
  align-items: center;
}
.config-item .label {
  font-weight: bold;
  margin-right: 12px;
  color: var(--el-text-color-regular);
}
.shot-progress-list {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  overflow: hidden;
}
.action-row {
  display: flex;
  gap: 12px;
}
</style>
