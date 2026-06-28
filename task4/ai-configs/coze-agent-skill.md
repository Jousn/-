# Coze Agent Skill 配置

## Agent 名称：动物世界

## Skill 配置

```markdown
# 动物世界 Agent Skill

## 功能描述
用户输入一个动物相关的主题或描述，自动生成一段10秒的视频（包含音频和视频）。

## 核心能力
1. **智能提示词增强**：通过大模型丰富用户输入，生成更详细的视频描述
2. **视频生成**：调用视频生成工作流，创建10秒动物主题视频
3. **音频合成**：生成匹配视频内容的音频
4. **智能缓存**：检查数据库避免重复生成相同内容
5. **视频拼接**：自动合成视频+音频的完整作品

## 使用限制
1. **主题限制**：仅接受动物相关的内容输入
2. **输出提醒**：对于非动物主题的输入，会提示用户"请输入与动物相关的内容"
3. **生成时长**：固定生成10秒视频
4. **响应时间**：视频生成+音频合成约需30-60秒

## 工作流程
1. 接收用户输入 → 验证主题相关性
2. 大模型提示词增强 → 生成详细视频描述
3. 数据库查询 → 检查是否已生成过类似视频
4. 视频生成工作流 → 调用视频生成API
5. 音频生成工作流 → 调用音频合成API
6. 视频拼接 → 合成视频+音频
7. 数据库存储 → 缓存生成结果
8. 输出视频地址 → 返回给用户

## 示例输入
- "一只可爱的小猫咪在花园里玩耍"
- "大象在河边喝水"
- "森林里的松鼠在找坚果"

## 示例输出
- 视频地址：https://cdn.example.com/animal_video_001.mp4
- 描述：10秒视频，包含画面和配音
```

---

## 2. 工作流配置

### 主工作流（主路径）

```yaml
name: animal_video_main_workflow
description: 动物视频生成主流程
version: 1.0

steps:
  # 步骤1：提示词增强
  - name: enhance_prompt
    type: llm_call
    model: gpt-4
    input:
      user_input: "{{user_input}}"
    prompt_template: |
      用户输入：{{user_input}}
      
      请将用户的简单描述增强为详细的视频生成提示词，包含：
      1. 动物外貌描述（颜色、体型、特征）
      2. 动物行为描述（动作、表情、姿态）
      3. 环境描述（背景、光照、氛围）
      4. 镜头描述（角度、景别、运动）
      5. 音效描述（动物声音、环境音效）
      
      输出格式：JSON
      {
        "enhanced_prompt": "详细的视频描述",
        "animal_features": {...},
        "environment": {...},
        "camera": {...},
        "audio": {...}
      }
    output_mapping:
      enhanced_prompt: "{{output.enhanced_prompt}}"
      animal_features: "{{output.animal_features}}"

  # 步骤2：数据库查询
  - name: check_cache
    type: database_query
    database: animal_video_cache
    query: |
      SELECT video_url, audio_url FROM video_cache
      WHERE prompt_similarity > 0.8
      AND animal_type = '{{animal_features.type}}'
      LIMIT 1
    condition:
      - if: "{{query_result.length > 0}}"
        then: return_cached_result
      - else: continue_to_generation

  # 步骤3：视频生成
  - name: generate_video
    type: workflow_call
    workflow: video_generation_workflow
    input:
      prompt: "{{enhanced_prompt}}"
      duration: 10
      resolution: "1080p"
    output_mapping:
      video_url: "{{output.video_url}}"

  # 步骤4：音频生成
  - name: generate_audio
    type: workflow_call
    workflow: audio_generation_workflow
    input:
      animal_type: "{{animal_features.type}}"
      behavior: "{{animal_features.behavior}}"
      environment: "{{environment.type}}"
    output_mapping:
      audio_url: "{{output.audio_url}}"

  # 步骤5：视频拼接
  - name: merge_video_audio
    type: api_call
    service: video_merge_service
    input:
      video_url: "{{video_url}}"
      audio_url: "{{audio_url}}"
    output_mapping:
      final_video_url: "{{output.final_video_url}}"

  # 步骤6：存储到数据库
  - name: save_to_cache
    type: database_insert
    database: animal_video_cache
    data:
      prompt: "{{user_input}}"
      enhanced_prompt: "{{enhanced_prompt}}"
      animal_type: "{{animal_features.type}}"
      video_url: "{{final_video_url}}"
      created_at: "{{timestamp}}"

  # 步骤7：返回结果
  - name: return_result
    type: output
    data:
      video_url: "{{final_video_url}}"
      description: "10秒动物主题视频，已包含配音"
      generation_time: "{{execution_time}}"
```

---

### 视频生成工作流

```yaml
name: video_generation_workflow
description: 生成10秒动物主题视频
version: 1.0

steps:
  - name: call_video_api
    type: api_call
    service: video_generation_api
    input:
      prompt: "{{prompt}}"
      duration: 10
      style: "realistic"
      resolution: "1080p"
    timeout: 45
    retry: 3
    output_mapping:
      video_url: "{{output.video_url}}"
      thumbnail: "{{output.thumbnail}}"

  - name: quality_check
    type: llm_call
    model: gpt-4-vision
    input:
      video_url: "{{video_url}}"
    prompt_template: |
      检查生成的视频是否符合以下要求：
      1. 视频时长是否为10秒
      2. 动物是否清晰可见
      3. 背景是否匹配描述
      4. 动作是否流畅自然
      
      输出：{"quality_score": 0-10, "issues": []}
    condition:
      - if: "{{output.quality_score < 7}}"
        then: regenerate_video
      - else: return_video

  - name: return_video
    type: output
    data:
      video_url: "{{video_url}}"
      quality_score: "{{quality_check.quality_score}}"
```

---

### 音频生成工作流

```yaml
name: audio_generation_workflow
description: 生成动物配音和环境音效
version: 1.0

steps:
  - name: generate_animal_sound
    type: api_call
    service: audio_generation_api
    input:
      animal_type: "{{animal_type}}"
      behavior: "{{behavior}}"
      duration: 10
    output_mapping:
      animal_audio_url: "{{output.audio_url}}"

  - name: generate_environment_sound
    type: api_call
    service: audio_generation_api
    input:
      environment: "{{environment}}"
      duration: 10
    output_mapping:
      environment_audio_url: "{{output.audio_url}}"

  - name: mix_audio
    type: api_call
    service: audio_mix_service
    input:
      animal_audio: "{{animal_audio_url}}"
      environment_audio: "{{environment_audio_url}}"
      animal_volume: 0.8
      environment_volume: 0.3
    output_mapping:
      final_audio_url: "{{output.final_audio_url}}"

  - name: return_audio
    type: output
    data:
      audio_url: "{{final_audio_url}}"
```

---

## 3. 数据库设计

### video_cache 表结构

```sql
CREATE TABLE animal_video_cache (
    id UUID PRIMARY KEY,
    user_input TEXT NOT NULL,
    enhanced_prompt TEXT NOT NULL,
    animal_type VARCHAR(50),
    video_url TEXT NOT NULL,
    audio_url TEXT,
    final_video_url TEXT NOT NULL,
    prompt_embedding VECTOR(1536),  -- 用于相似度匹配
    created_at TIMESTAMP DEFAULT NOW(),
    quality_score INT,
    generation_time INT  -- 生成耗时（秒）
);

CREATE INDEX idx_animal_type ON animal_video_cache(animal_type);
CREATE INDEX idx_prompt_embedding ON animal_video_cache USING ivfflat(prompt_embedding);
```

---

## 4. 配置解决的问题

### 问题1：用户输入过于简单
- **问题描述**：用户输入"一只猫"，无法直接生成高质量视频
- **解决方案**：通过大模型提示词增强，生成详细描述
- **场景生效**：所有用户输入都会经过增强处理
- **迭代优化**：
  - 第一版：仅增加动物描述
  - 第二版：增加环境、镜头、音效描述
  - 第三版：结构化输出（JSON格式），便于后续处理

### 问题2：重复生成浪费资源
- **问题描述**：多个用户输入相似内容，重复生成相同视频
- **解决方案**：数据库缓存 + 相似度匹配（向量检索）
- **场景生效**：每次生成前先查询缓存
- **迭代优化**：
  - 第一版：关键词匹配（精确匹配）
  - 第二版：向量相似度匹配（threshold: 0.8）
  - 第三版：增加动物类型分类索引，提升查询效率

### 问题3：视频和音频不匹配
- **问题描述**：独立生成的视频和音频可能不协调
- **解决方案**：统一提示词增强，确保视频和音频描述一致
- **场景生效**：视频生成和音频生成使用相同的增强提示词
- **迭代优化**：
  - 第一版：独立生成，人工拼接
  - 第二版：统一增强，自动拼接
  - 第三版：增加音频混音参数（动物音量0.8，环境音量0.3）



