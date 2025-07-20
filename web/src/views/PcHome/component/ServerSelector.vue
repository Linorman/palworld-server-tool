<template>
  <div class="server-selector">
    <n-space align="center" justify="space-between">
      <n-space align="center">
        <n-icon size="18" color="#4098fc">
          <server-outlined />
        </n-icon>
        <n-select
          v-model:value="selectedServerId"
          :options="serverOptions"
          :loading="loading"
          placeholder="选择服务器"
          :render-label="renderLabel"
          :render-tag="renderTag"
          size="small"
          style="min-width: 200px"
          @update:value="onServerChange"
        />
        <n-button
          v-if="isAdmin"
          size="small"
          type="primary"
          ghost
          @click="showManageModal = true"
        >
          <template #icon>
            <n-icon>
              <settings-outlined />
            </n-icon>
          </template>
          管理服务器
        </n-button>
      </n-space>
      
      <server-status-indicator @select-server="onServerSelect" />
    </n-space>

    <!-- 服务器管理弹窗 -->
    <n-modal
      v-model:show="showManageModal"
      preset="dialog"
      title="服务器管理"
      style="width: 80%; max-width: 1200px"
    >
      <server-management @close="showManageModal = false" />
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, h } from 'vue';
import { NSelect, NButton, NModal, NSpace, NIcon, NTag } from 'naive-ui';
import { DeploymentUnitOutlined, SettingOutlined } from '@vicons/antd';
import serverStore from '@/stores/model/server';
import userStore from '@/stores/model/user';
import ApiService from '@/service/api';
import { useMessage } from 'naive-ui';
import ServerManagement from './ServerManagement.vue';
import ServerStatusIndicator from './ServerStatusIndicator.vue';

const message = useMessage();
const apiService = new ApiService();

const selectedServerId = ref(null);
const showManageModal = ref(false);
const loading = ref(false);

const servers = computed(() => serverStore().getServers());
const currentServer = computed(() => serverStore().getCurrentServer());
const isAdmin = computed(() => userStore().getLoginInfo().isLogin);

// 服务器选项
const serverOptions = computed(() => {
  return servers.value
    .filter(server => server.enabled)
    .map(server => ({
      label: server.name,
      value: server.id,
      server: server
    }));
});

// 自定义渲染标签
const renderLabel = (option) => {
  return h('div', { style: 'display: flex; align-items: center; gap: 8px' }, [
    h(NTag, { 
      size: 'small', 
      type: option.server.enabled ? 'success' : 'default',
      bordered: false 
    }, { default: () => option.server.enabled ? '在线' : '离线' }),
    h('span', option.label),
    option.server.description && h('span', { style: 'color: #999; font-size: 12px' }, ` - ${option.server.description}`)
  ]);
};

// 自定义渲染选中标签
const renderTag = ({ option }) => {
  return h('div', { style: 'display: flex; align-items: center; gap: 4px' }, [
    h(NTag, { 
      size: 'small', 
      type: option.server.enabled ? 'success' : 'default',
      bordered: false 
    }, { default: () => option.server.enabled ? '●' : '○' }),
    h('span', option.label)
  ]);
};

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

// 服务器切换
const onServerChange = (serverId) => {
  const server = servers.value.find(s => s.id === serverId);
  if (server) {
    serverStore().setCurrentServer(server);
    message.success(`已切换到服务器：${server.name}`);
  }
};

// 从状态指示器选择服务器
const onServerSelect = (server) => {
  selectedServerId.value = server.id;
  serverStore().setCurrentServer(server);
  message.success(`已切换到服务器：${server.name}`);
};

// 监听当前服务器变化
watch(currentServer, (newServer) => {
  if (newServer) {
    selectedServerId.value = newServer.id;
  }
});

// 初始化
onMounted(async () => {
  await loadServers();
  // 设置初始选中的服务器
  if (currentServer.value) {
    selectedServerId.value = currentServer.value.id;
  }
});
</script>

<style scoped>
.server-selector {
  padding: 8px 0;
}
</style> 