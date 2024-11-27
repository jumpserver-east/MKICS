<template>
    <div class="kf-page">
      <!-- 表格组件 -->
      <o-table ref="tableRef" :table-config="tableConfig">
        <!-- 渲染 slotColumns -->
        <template v-for="(item, index) in slotColumns" :key="index" #[item.slot!]="{ row }">
          <!-- 自定义渲染逻辑 -->
          <span v-if="item.prop === 'staff_list'">
            {{ row.staff_list?.join('; ') || '-' }}
          </span>
          <span v-else-if="item.prop === 'staffs'">
            {{ (row.staffs as Array<IStaff>)?.map((staff) => staff.staffname).join('; ') || '-' }}
          </span>
          <span v-else-if="item.prop === 'transfer_keywords'">
            {{ row[item.prop]?.split(',').join('; ') || '-' }}
          </span>
          <span v-else>
            {{ row[item.prop] || '-' }}
          </span>
        </template>
  
        <!-- 表格顶部按钮 -->
        <template #table-top>
          <el-button type="primary" @click="handleAdd">添加客服</el-button>
        </template>
      </o-table>
    </div>
  </template>
  
  <script lang="ts" setup>
  import { ref } from 'vue';
  import router from '@/router';
  import { useuuidDelete } from '@/hooks';
  import { Api } from '@/api/common/enum';
  import { deleteKFApi } from '@/api/kf';
  import type { ITableConfig, TableInstance } from '@/types';
  import type { IStaff } from '@/api/staff/model';
  import type { IKF } from '@/api/kf/model';
  
  // 配置表格信息
  const tableConfig: ITableConfig = {
    api: `${Api.kf}`, // 使用 kf 的接口路径
    headers: [
      { prop: 'kfname', label: '客服名称' },
      { prop: 'kfid', label: '客服 ID' },
      { prop: 'kfplatform', label: '客服平台' },
      { prop: 'botid', label: '机器人 ID' },
      { prop: 'botplatform', label: '机器人平台' },
      { prop: 'status', label: '状态' },
      { prop: 'receive_priority', label: '接待优先级' },
      { prop: 'receive_rule', label: '接待规则' },
      { prop: 'chat_timeout', label: '会话超时（秒）' },
      { prop: 'bot_timeout', label: '机器人超时（秒）' },
      { prop: 'bot_timeout_msg', label: '机器人超时提示' },
      { prop: 'bot_welcome_msg', label: '机器人欢迎语' },
      { prop: 'staff_welcome_msg', label: '接待人员欢迎语' },
      { prop: 'unmanned_msg', label: '无人接待提示' },
      { prop: 'chatend_msg', label: '会话结束语' },
      { prop: 'transfer_keywords', label: '转接关键词', slot: 'transfer_keywords' },
      { prop: 'staffs', label: '人员名称列表', slot: 'staffs' },
    ],
    operations: {
      width: 150, // 操作列宽度
      buttons: [
        {
          text: '编辑',
          show: true,
          click: ({ row }) => {
            router.push(`/kf/edit/${row.uuid}`);
          },
        },
        {
          text: '删除',
          type: 'danger',
          show: true,
          click: ({ row }) => {
            handleDelete(row.uuid);
          },
        },
      ],
    },
  };
  
  // 筛选带有 slot 的字段
  const slotColumns = tableConfig.headers.filter((header) => header.slot);
  
  // 表格实例
  const tableRef = ref<TableInstance>();
  
  // 添加客服
  const handleAdd = () => {
    router.push('/kf/add');
  };
  
  // 删除客服
  const handleDelete = (uuid: string) => {
    const { onDelete } = useuuidDelete();
    onDelete(deleteKFApi, uuid, tableRef); // 使用 kf 的删除接口
  };
  </script>
  