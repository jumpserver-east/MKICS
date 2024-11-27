<template>
    <div class="receptionist-details">
      <!-- 接待人员列表 -->
      <o-table ref="tableRef" :table-config="tableConfig">
        <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
          <span v-if="item.prop === 'status'">
            <!-- 根据 status 显示不同状态 -->
            <span :class="{'status-active': row[item.prop] === 0, 'status-disabled': row[item.prop] === 1}">
              {{ row[item.prop] === 0 ? '接待中' : '停止接待' }}
            </span>
          </span>
          <span v-else>
            <!-- 处理接待人员字段 -->
            <template v-if="item.prop === 'userid'">
              {{ row[item.prop] }}
            </template>
          </span>
        </template>

        <!-- 表格顶部操作：添加接待人员 -->
        <template #table-top>
          <el-input
            v-model="newReceptionists"
            placeholder="输入接待人员ID, 用逗号隔开"
            @keyup.enter="handleAddReceptionists"
          />
          <el-button type="primary" @click="handleAddReceptionists">添加接待人员</el-button>
          <el-button type="danger" @click="handleDeleteSelected">批量删除接待人员</el-button>
        </template>
      </o-table>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getReceptionistListApi, deleteReceptionistApi, createReceptionistApi } from '@/api/wecom/receptionist'
import { Api } from '@/api/common/enum'
import type { ITableConfig, TableInstance } from '@/types'

const tableRef = ref<TableInstance>()
const route = useRoute()
const router = useRouter()

// 获取 kfid 参数
const kfid = route.params.kfid as string

// 表格配置
const tableConfig: ITableConfig = {
  api: `${Api.wecomreceptionist}/${kfid}`,
  headers: [
    { prop: 'userid', label: '接待人员ID' },
    { prop: 'status', label: '状态', slot: 'status' }
  ],
  operations: {
    width: 120,
    buttons: [
      {
        text: '删除',
        type: 'danger',
        show: true,
        click: ({ row }) => handleDelete(row.userid)  // 删除单个接待人员
      }
    ]
  }
}

const slotColumns = tableConfig.headers.filter((header) => header.slot)

// 加载接待人员列表
const receptionistList = ref<any[]>([])

// 输入的接待人员ID列表
const newReceptionists = ref('')

// 新增的接待人员ID列表
const handleAddReceptionists = async () => {
  const userIds = newReceptionists.value.split(',').map(id => id.trim())

  if (userIds.length > 0) {
    try {
      // 发送一个包含多个 userid 的请求
      await createReceptionistApi(kfid, {
        userid_list: userIds,
        open_kfid: kfid
      })

      // 清空输入框
      newReceptionists.value = ''
    } catch (error) {
      console.error('添加接待人员失败', error)
    }
  } else {
    console.error('接待人员ID不能为空')
  }
}

// 删除接待人员
const handleDelete = async (userid: string) => {
  try {
    await deleteReceptionistApi(kfid, {
      userid_list: [userid], // 确保传递一个数组
      open_kfid: kfid
    })
  } catch (error) {
    console.error('删除接待人员失败', error)
  }
}

// 批量删除选中的接待人员
const handleDeleteSelected = async () => {
  const selectedUserIds = receptionistList.value.filter(item => item.selected).map(item => item.userid)
  if (selectedUserIds.length > 0) {
    try {
      // 批量删除选中的接待人员
      await deleteReceptionistApi(kfid, {
        userid_list: selectedUserIds, // 不再嵌套数组
        open_kfid: kfid
      })
    } catch (error) {
      console.error('批量删除接待人员失败', error)
    }
  } else {
    console.error('请先选择接待人员进行删除')
  }
}
</script>