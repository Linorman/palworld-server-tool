<template>
  <div class="server-status-indicator">
    <n-space align="center">
      <n-icon size="16" color="#4098fc">
        <server-outlined />
      </n-icon>
      <n-text style="font-size: 12px; color: #999">服务器状态</n-text>
      <n-space>
        <n-tooltip v-for="server in servers" :key="server.id" placement="top">
          <template #trigger>
            <n-tag
              :type="getServerStatusType(server)"
              size="small"
              round
              :bordered="false"
              style="cursor: pointer"
              @click="$emit('select-server', server)"
            >
              {{ server.name }}
            </n-tag>
          </template>
          <div>
            <p><strong>{{ server.name }}</strong></p>
            <p>ID: {{ server.id }}</p>
            <p>状态: {{ server.enabled ? '在线' : '离线' }}</p>
            <p v-if="server.description">{{ server.description }}</p>
            <p v-if="serverMetrics[server.id]">
              玩家: {{ serverMetrics[server.id].online_players || 0 }}/{{ serverMetrics[server.id].max_players || 0 }}
            </p>
            <p v-if="serverMetrics[server.id]">
              FPS: {{ serverMetrics[server.id].server_fps || 0 }}
            </p>
          </div>
        </n-tooltip>
      </n-space>
    </n-space>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { NSpace, NIcon, NTag, NTooltip, NText } from 'naive-ui';
import { ServerOutlined } from '@vicons/antd';
import serverStore from '@/stores/model/server';
import ApiService from '@/service/api';

const emit = defineEmits(['select-server']);

const serverMetrics = ref({});
const apiService = new ApiService();

const servers = computed(() => serverStore().getServers());

// 获取服务器状态类型
const getServerStatusType = (server) => {
  if (!server.enabled) return 'default';
  
  const metrics = serverMetrics.value[server.id];
  if (!metrics) return 'warning';
  
  // 根据FPS判断服务器状态
  const fps = metrics.server_fps || 0;
  if (fps >= 30) return 'success';
  if (fps >= 15) return 'warning';
  return 'error';
};

// 获取所有服务器的状态
const loadServerMetrics = async () => {
  const promises = servers.value.map(async (server) => {
    if (!server.enabled) return;
    
    try {
      const { data } = await apiService.getServerMetrics(server.id);
      if (data.value) {
        serverMetrics.value[server.id] = data.value;
      }
    } catch (error) {
      console.error(`Failed to load metrics for server ${server.id}:`, error);
    }
  });
  
  await Promise.all(promises);
};

// 定时更新服务器状态
let statusInterval = null;

onMounted(() => {
  loadServerMetrics();
  statusInterval = setInterval(loadServerMetrics, 30000); // 每30秒更新一次
});

onUnmounted(() => {
  if (statusInterval) {
    clearInterval(statusInterval);
  }
});
</script>

<style scoped>
.server-status-indicator {
  padding: 4px 0;
}
</style> 