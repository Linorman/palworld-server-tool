<template>
  <div class="server-management">
    <n-space vertical>
      <!-- 头部操作 -->
      <n-space justify="space-between">
        <n-h3>服务器管理</n-h3>
        <n-button type="primary" @click="showAddModal = true">
          <template #icon>
            <n-icon>
              <plus-outlined />
            </n-icon>
          </template>
          添加服务器
        </n-button>
      </n-space>

      <!-- 服务器列表 -->
      <n-data-table
        :columns="columns"
        :data="servers"
        :loading="loading"
        :pagination="false"
        :bordered="false"
        size="small"
      />
    </n-space>

    <!-- 添加/编辑服务器弹窗 -->
    <n-modal
      v-model:show="showAddModal"
      preset="dialog"
      title="添加服务器"
      style="width: 600px"
    >
      <server-form
        :server="editingServer"
        @submit="handleSubmit"
        @cancel="handleCancel"
      />
    </n-modal>

    <n-modal
      v-model:show="showEditModal"
      preset="dialog"
      title="编辑服务器"
      style="width: 600px"
    >
      <server-form
        :server="editingServer"
        @submit="handleSubmit"
        @cancel="handleCancel"
      />
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue';
import { NDataTable, NH3, NButton, NModal, NSpace, NIcon, NTag, NPopconfirm } from 'naive-ui';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@vicons/antd';
import serverStore from '@/stores/model/server';
import ApiService from '@/service/api';
import { useMessage } from 'naive-ui';
import ServerForm from './ServerForm.vue';

const emit = defineEmits(['close']);

const message = useMessage();
const apiService = new ApiService();

const loading = ref(false);
const showAddModal = ref(false);
const showEditModal = ref(false);
const editingServer = ref(null);

const servers = computed(() => serverStore().getServers());

// 表格列定义
const columns = [
  {
    title: '服务器名称',
    key: 'name',
    render: (row) => {
      return h('div', { style: 'display: flex; align-items: center; gap: 8px' }, [
        h(NTag, { 
          size: 'small', 
          type: row.enabled ? 'success' : 'default',
          bordered: false 
        }, { default: () => row.enabled ? '在线' : '离线' }),
        h('span', row.name)
      ]);
    }
  },
  {
    title: '服务器ID',
    key: 'id',
    width: 120
  },
  {
    title: '描述',
    key: 'description',
    ellipsis: true
  },
  {
    title: 'RCON地址',
    key: 'rcon',
    render: (row) => row.rcon?.address || '-'
  },
  {
    title: 'REST地址',
    key: 'rest',
    render: (row) => row.rest?.address || '-'
  },
  {
    title: '状态',
    key: 'enabled',
    width: 80,
    render: (row) => {
      return h(NTag, { 
        type: row.enabled ? 'success' : 'default',
        size: 'small'
      }, { default: () => row.enabled ? '启用' : '禁用' });
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row) => {
      return h('div', { style: 'display: flex; gap: 8px' }, [
        h(NButton, {
          size: 'small',
          type: 'primary',
          ghost: true,
          onClick: () => handleEdit(row)
        }, { default: () => '编辑' }),
        h(NPopconfirm, {
          onPositiveClick: () => handleDelete(row.id)
        }, {
          default: () => '确定删除这个服务器吗？',
          trigger: () => h(NButton, {
            size: 'small',
            type: 'error',
            ghost: true
          }, { default: () => '删除' })
        })
      ]);
    }
  }
];

// 加载服务器列表
const loadServers = async () => {
  try {
    loading.value = true;
    const { data } = await apiService.getServers();
    if (data.value && data.value.servers) {
      serverStore().setServers(data.value.servers);
    }
  } catch (error) {
    message.error('加载服务器列表失败');
  } finally {
    loading.value = false;
  }
};

// 编辑服务器
const handleEdit = (server) => {
  editingServer.value = { ...server };
  showEditModal.value = true;
};

// 删除服务器
const handleDelete = async (serverId) => {
  try {
    await apiService.deleteServer(serverId);
    serverStore().removeServer(serverId);
    message.success('服务器删除成功');
  } catch (error) {
    message.error('服务器删除失败');
  }
};

// 表单提交
const handleSubmit = async (formData) => {
  try {
    if (editingServer.value?.id) {
      // 更新服务器
      await apiService.updateServer(editingServer.value.id, formData);
      serverStore().updateServer(editingServer.value.id, formData);
      message.success('服务器更新成功');
    } else {
      // 创建服务器
      const { data } = await apiService.createServer(formData);
      if (data.value) {
        serverStore().addServer(data.value);
        message.success('服务器创建成功');
      }
    }
    handleCancel();
  } catch (error) {
    message.error(editingServer.value?.id ? '服务器更新失败' : '服务器创建失败');
  }
};

// 取消操作
const handleCancel = () => {
  showAddModal.value = false;
  showEditModal.value = false;
  editingServer.value = null;
};

// 初始化
onMounted(() => {
  loadServers();
});
</script>

<style scoped>
.server-management {
  padding: 16px;
}
</style> 