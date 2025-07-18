import { defineStore } from "pinia";
import { ref } from "vue";

const serverStore = defineStore(
  "server",
  () => {
    const servers = ref([]);
    const currentServer = ref(null);
    const isLoading = ref(false);

    const setServers = (serverList) => {
      servers.value = serverList;
      // 如果没有选中的服务器，自动选择第一个启用的服务器
      if (!currentServer.value && serverList.length > 0) {
        const firstEnabledServer = serverList.find(server => server.enabled);
        if (firstEnabledServer) {
          currentServer.value = firstEnabledServer;
        }
      }
    };

    const setCurrentServer = (server) => {
      currentServer.value = server;
    };

    const getCurrentServer = () => {
      return currentServer.value;
    };

    const getServers = () => {
      return servers.value;
    };

    const getEnabledServers = () => {
      return servers.value.filter(server => server.enabled);
    };

    const setLoading = (loading) => {
      isLoading.value = loading;
    };

    const getLoading = () => {
      return isLoading.value;
    };

    const addServer = (server) => {
      servers.value.push(server);
    };

    const updateServer = (serverId, serverData) => {
      const index = servers.value.findIndex(s => s.id === serverId);
      if (index !== -1) {
        servers.value[index] = { ...servers.value[index], ...serverData };
        // 如果更新的是当前选中的服务器，也要更新当前服务器
        if (currentServer.value && currentServer.value.id === serverId) {
          currentServer.value = { ...currentServer.value, ...serverData };
        }
      }
    };

    const removeServer = (serverId) => {
      servers.value = servers.value.filter(s => s.id !== serverId);
      // 如果删除的是当前选中的服务器，选择第一个启用的服务器
      if (currentServer.value && currentServer.value.id === serverId) {
        const firstEnabledServer = servers.value.find(server => server.enabled);
        currentServer.value = firstEnabledServer || null;
      }
    };

    const getServerById = (serverId) => {
      return servers.value.find(s => s.id === serverId);
    };

    return {
      servers,
      currentServer,
      isLoading,
      setServers,
      setCurrentServer,
      getCurrentServer,
      getServers,
      getEnabledServers,
      setLoading,
      getLoading,
      addServer,
      updateServer,
      removeServer,
      getServerById,
    };
  },
  {
    persist: {
      key: "server-store",
      storage: localStorage,
      paths: ["currentServer"] // 只持久化当前选中的服务器
    },
  }
);

export default serverStore; 