<template>
    <div class="maxkb-config">
      <el-card>
        <div class="header">
          <h2>MaxKB 配置</h2>
          <el-button type="primary" @click="editMode = !editMode">
            {{ editMode ? "取消修改" : "编辑配置" }}
          </el-button>
        </div>
  
        <el-form
          :model="formData"
          :rules="formRules"
          ref="configForm"
          label-position="top"
          :disabled="!editMode"
        >
          <el-form-item label="Base URL：" prop="base_url">
            <el-input v-model="formData.base_url" placeholder="请输入 Base URL" />
          </el-form-item>
  
          <el-form-item label="API Key：" prop="api_key">
            <el-input
              v-model="formData.api_key"
              placeholder="请输入 API Key"
              show-password
            />
          </el-form-item>
  
          <el-form-item v-if="editMode">
            <el-button type="primary" @click="submitConfig">保存配置</el-button>
            <el-button @click="cancelEdit">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </template>
  
  <script lang="ts" setup>
  import { reactive, ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getConfigApi, updateConfigApi } from '@/api/maxkb'
  
  interface MaxkbConf {
    base_url: string
    api_key: string
  }
  
  const formData = reactive<MaxkbConf>({
    base_url: '',
    api_key: ''
  })
  
  const formRules = {
    base_url: [{ required: true, message: '请输入 Base URL', trigger: 'blur' }],
    api_key: [{ required: true, message: '请输入 API Key', trigger: 'blur' }]
  }
  
  const editMode = ref(false)
  const configForm = ref()
  
  // 获取配置数据
  const fetchConfig = async () => {
    try {
      const { data } = await getConfigApi()
      Object.assign(formData, data)
    } catch (error) {
      ElMessage.error('获取配置失败，请稍后重试')
    }
  }
  
  // 提交配置
  const submitConfig = () => {
    configForm.value.validate(async (valid: boolean) => {
      if (valid) {
        try {
          await updateConfigApi(formData)
          ElMessage.success('配置更新成功')
          editMode.value = false
          fetchConfig()
        } catch (error) {
          ElMessage.error('更新配置失败，请稍后重试')
        }
      }
    })
  }
  
  // 取消编辑
  const cancelEdit = () => {
    editMode.value = false
    fetchConfig()
  }
  
  onMounted(() => {
    fetchConfig()
  })
  </script>
  
  <style scoped>
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  </style>
  