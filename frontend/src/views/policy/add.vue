<template>
  <o-form-wrap title="策略管理" @confirm="onConfirm">
    <el-form
      ref="ruleForm"
      class="w-100"
      :model="formData"
      :rules="formRules"
      label-position="top"
    >
      <el-form-item label="策略名称：" prop="policyname">
        <el-input v-model="formData.policyname" placeholder="请输入策略名称" />
      </el-form-item>

      <el-form-item label="最大接待数量：" prop="max_count">
        <el-input
          v-model="formData.max_count"
          type="number"
          placeholder="请输入最大接待数量"
        />
      </el-form-item>

      <el-form-item label="重复策略：" prop="repeat">
        <el-input
          v-model="formData.repeat"
          type="number"
          placeholder="请输入重复策略"
        />
      </el-form-item>

      <el-form-item label="工作日：" prop="week">
        <el-input v-model="formData.week" placeholder="请输入工作日标记（如：1111100）" />
      </el-form-item>

      <el-form-item label="工作时间：" prop="work_times">
        <el-button type="primary" @click="addWorkTime">添加工作时间</el-button>
        <div v-for="(time, index) in formData.work_times" :key="index" class="work-time">
          <el-input v-model="time.start_time" placeholder="开始时间"></el-input>
          <el-input v-model="time.end_time" placeholder="结束时间"></el-input>
          <el-button @click="removeWorkTime(index)" type="danger">删除</el-button>
        </div>
      </el-form-item>
    </el-form>
  </o-form-wrap>
</template>

<script lang="ts" setup>
import { onBeforeMount, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'

import { createPolicyApi, getPolicyApi, updatePolicyApi } from '@/api/policy'
import type { IPolicy } from '@/api/policy/model'

import type { TLoading } from '@/types'
import type { FormInstance, FormRules } from 'element-plus'

const route = useRoute()
const router = useRouter()

const isEditing = ref(false)
const uuid = route.params.uuid as string;

isEditing.value = uuid ? true : false

const ruleForm = ref<FormInstance>()
const formData = reactive<IPolicy>({
  uuid:'',
  policyname: '',
  max_count: 100,
  repeat: 1,
  week: '',
  work_times: [{ start_time: '', end_time: '' }]
})

const formRules = reactive<FormRules>({
  policyname: [{ required: true, message: '请输入策略名称', trigger: 'blur' }],
  max_count: [{ required: true, message: '请输入最大接待数量', trigger: 'blur' }],
  repeat: [{ required: true, message: '请输入重复策略', trigger: 'blur' }],
  week: [{ required: true, message: '请输入工作日标记', trigger: 'blur' }],
  work_times: [{ required: true, type: 'array', min: 1, message: '请添加至少一个工作时间', trigger: 'blur' }]
})

const policyInfo = async () => {
  const { data } = await getPolicyApi(uuid);
  Object.assign(formData, data)
}

const onConfirm = (loading: TLoading) => {
  ruleForm.value?.validate(async (valid) => {
    if (valid) {
      loading(true)
      if (isEditing.value) {
        await updatePolicyApi(formData)
        ElMessage.success('编辑成功')
      } else {
        await createPolicyApi(formData)
        ElMessage.success('添加成功')
      }

      loading(false)
      router.back()
    } else {
      return
    }
  })
}

const addWorkTime = () => {
  if (Array.isArray(formData.work_times)) {
    formData.work_times.push({ start_time: '', end_time: '' });
  }
}

const removeWorkTime = async (index: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该工作时间吗？',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    // 用户点击确认后执行删除操作
    formData.work_times?.splice(index, 1)
    ElMessage.success('工作时间已删除')
  } catch {
  }
}

onBeforeMount(() => {
  if (isEditing.value) policyInfo()
})
</script>
