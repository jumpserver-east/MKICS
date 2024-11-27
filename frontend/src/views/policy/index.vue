<template>
  <o-table ref="tableRef" :table-config="tableConfig">
    <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
      <!-- 处理 work_times 字段 -->
      <span v-if="item.prop === 'work_times'">
        <span v-for="(time, idx) in row[item.prop]" :key="idx">
          {{ time.start_time }} - {{ time.end_time }}
          <span v-if="idx !== row[item.prop].length - 1">, </span>
        </span>
      </span>
      <!-- 处理其他字段 -->
      <span v-else>
        {{ row[item.prop] ? row[item.prop].name : '-' }}
      </span>
    </template>

    <template #table-top>
      <el-button type="primary" @click="handleAdd">添加</el-button>
    </template>
  </o-table>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import router from '@/router'
import { useuuidDelete } from '@/hooks'
import { Api } from '@/api/common/enum'
import { deletePolicyApi } from '@/api/policy'
import type { ITableConfig, TableInstance } from '@/types'
const tableConfig: ITableConfig = {
  api: `${Api.policy}`,
  headers: [
    {
      prop: 'policyname',
      label: '策略名称'
    },
    {
      prop: 'max_count',
      label: '最大接待数量'
    },
    {
      prop: 'repeat',
      label: '重复策略'
    },
    {
      prop: 'week',
      label: '自定义工作日'
    },
    {
      prop: 'work_times',
      label: '工作时间',
      slot: 'work_times'
    }
  ],
  operations: {
    width: 115,
    buttons: [
      {
        text: '编辑',
        show: true,
        click: ({ row }) => {
          router.push(`/policy/edit/${row.uuid}`)
        }
      },
      {
        text: '删除',
        type: 'danger',
        show: true,
        click: ({ row }) => {
          handleDelete(row.uuid)
        }
      }
    ]
  }
}

const slotColumns = tableConfig.headers.filter((header) => header.slot)

const tableRef = ref<TableInstance>()

const handleAdd = () => {
  router.push('/policy/add')
}

const handleDelete = (uuid: string) => {
  const { onDelete } = useuuidDelete()
  onDelete(deletePolicyApi, uuid, tableRef)
}
</script>
