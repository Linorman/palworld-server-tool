import Service from "./service";

class ApiService extends Service {
  async login(param) {
    let data = param;
    return this.fetch(`/api/login`).post(data).json();
  }

  // 多服务器管理 API
  async getServers() {
    return this.fetch(`/api/servers`).get().json();
  }

  async getServerById(serverId) {
    return this.fetch(`/api/servers/${serverId}`).get().json();
  }

  async createServer(param) {
    let data = param;
    return this.fetch(`/api/servers`).post(data).json();
  }

  async updateServer(serverId, param) {
    let data = param;
    return this.fetch(`/api/servers/${serverId}`).put(data).json();
  }

  async deleteServer(serverId) {
    return this.fetch(`/api/servers/${serverId}`).delete().json();
  }

  // 多服务器操作 API
  async getServerInfo(serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/info`).get().json();
    }
    return this.fetch(`/api/server`).get().json();
  }

  async getServerMetrics(serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/metrics`).get().json();
    }
    return this.fetch(`/api/server/metrics`).get().json();
  }

  async getServerToolInfo() {
    return this.fetch(`/api/server/tool`).get().json();
  }

  async sendBroadcast(param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/broadcast`).post(data).json();
    }
    return this.fetch(`/api/server/broadcast`).post(data).json();
  }

  async shutdownServer(param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/shutdown`).post(data).json();
    }
    return this.fetch(`/api/server/shutdown`).post(data).json();
  }

  async getPlayerList(param, serverId = null) {
    const query = this.generateQuery(param);
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/players?${query}`).get().json();
    }
    return this.fetch(`/api/player?${query}`).get().json();
  }

  async getOnlinePlayerList(serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/online_players`).get().json();
    }
    return this.fetch(`/api/online_player`).get().json();
  }

  async getPlayer(param, serverId = null) {
    const { playerUid } = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/players/${playerUid}`).get().json();
    }
    return this.fetch(`/api/player/${playerUid}`).get().json();
  }

  async kickPlayer(param, serverId = null) {
    const { playerUid } = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/players/${playerUid}/kick`).post().json();
    }
    return this.fetch(`/api/player/${playerUid}/kick`).post().json();
  }

  async banPlayer(param, serverId = null) {
    const { playerUid } = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/players/${playerUid}/ban`).post().json();
    }
    return this.fetch(`/api/player/${playerUid}/ban`).post().json();
  }

  async unbanPlayer(param, serverId = null) {
    const { playerUid } = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/players/${playerUid}/unban`).post().json();
    }
    return this.fetch(`/api/player/${playerUid}/unban`).post().json();
  }

  async getGuildList(serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/guilds`).get().json();
    }
    return this.fetch(`/api/guild`).get().json();
  }

  async getGuild(param, serverId = null) {
    const { adminPlayerUid } = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/guilds/${adminPlayerUid}`).get().json();
    }
    return this.fetch(`/api/guild/${adminPlayerUid}`).get().json();
  }

  async getWhitelist(serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/whitelist`).get().json();
    }
    return this.fetch(`/api/whitelist`).get().json();
  }

  async addWhitelist(param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/whitelist`).post(data).json();
    }
    return this.fetch(`/api/whitelist`).post(data).json();
  }

  async removeWhitelist(param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/whitelist`).delete(data).json();
    }
    return this.fetch(`/api/whitelist`).delete(data).json();
  }

  async putWhitelist(param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/whitelist`).put(data).json();
    }
    return this.fetch(`/api/whitelist`).put(data).json();
  }

  async getRconCommands(serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/rcon`).get().json();
    }
    return this.fetch(`/api/rcon`).get().json();
  }

  async sendRconCommand(param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/rcon/send`).post(data).json();
    }
    return this.fetch(`/api/rcon/send`).post(data).json();
  }

  async addRconCommand(param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/rcon`).post(data).json();
    }
    return this.fetch(`/api/rcon`).post(data).json();
  }

  async putRconCommand(uuid, param, serverId = null) {
    let data = param;
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/rcon/${uuid}`).put(data).json();
    }
    return this.fetch(`/api/rcon/${uuid}`).put(data).json();
  }

  async removeRconCommand(uuid, serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/rcon/${uuid}`).delete().json();
    }
    return this.fetch(`/api/rcon/${uuid}`).delete().json();
  }

  async getBackupList(param, serverId = null) {
    const query = this.generateQuery(param);
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/backups?${query}`).get().json();
    }
    return this.fetch(`/api/backup?${query}`).get().json();
  }

  async removeBackup(uuid, serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/backups/${uuid}`).delete().json();
    }
    return this.fetch(`/api/backup/${uuid}`).delete().json();
  }

  async downloadBackup(uuid, serverId = null) {
    if (serverId) {
      return this.fetch(`/api/servers/${serverId}/backups/${uuid}`).get().blob();
    }
    return this.fetch(`/api/backup/${uuid}`).get().blob();
  }
}

export default ApiService;
