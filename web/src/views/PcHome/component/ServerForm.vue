<template>
  <div class="server-form">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-placement="left"
      label-width="120px"
      size="small"
    >
      <n-form-item label="服务器ID" path="id">
        <n-input 
          v-model:value="formData.id" 
          placeholder="请输入服务器唯一标识"
          :disabled="isEdit"
        />
      </n-form-item>

      <n-form-item label="服务器名称" path="name">
        <n-input 
          v-model:value="formData.name" 
          placeholder="请输入服务器名称"
        />
      </n-form-item>

      <n-form-item label="描述" path="description">
        <n-input 
          v-model:value="formData.description" 
          placeholder="请输入服务器描述（可选）"
          type="textarea"
          :rows="2"
        />
      </n-form-item>

      <n-form-item label="启用状态" path="enabled">
        <n-switch v-model:value="formData.enabled" />
      </n-form-item>

      <n-divider title-placement="left">RCON配置</n-divider>

      <n-form-item label="RCON地址" path="rcon.address">
        <n-input 
          v-model:value="formData.rcon.address" 
          placeholder="127.0.0.1:25575"
        />
      </n-form-item>

      <n-form-item label="RCON密码" path="rcon.password">
        <n-input 
          v-model:value="formData.rcon.password" 
          placeholder="请输入RCON密码"
          type="password"
          show-password-on="click"
        />
      </n-form-item>

      <n-form-item label="Base64编码" path="rcon.use_base64">
        <n-switch v-model:value="formData.rcon.use_base64" />
      </n-form-item>

      <n-form-item label="超时时间(秒)" path="rcon.timeout">
        <n-input-number 
          v-model:value="formData.rcon.timeout" 
          :min="1"
          :max="60"
          placeholder="5"
        />
      </n-form-item>

      <n-divider title-placement="left">REST API配置</n-divider>

      <n-form-item label="REST地址" path="rest.address">
        <n-input 
          v-model:value="formData.rest.address" 
          placeholder="http://127.0.0.1:8212"
        />
      </n-form-item>

      <n-form-item label="用户名" path="rest.username">
        <n-input 
          v-model:value="formData.rest.username" 
          placeholder="admin"
        />
      </n-form-item>

      <n-form-item label="密码" path="rest.password">
        <n-input 
          v-model:value="formData.rest.password" 
          placeholder="请输入REST API密码"
          type="password"
          show-password-on="click"
        />
      </n-form-item>

      <n-form-item label="超时时间(秒)" path="rest.timeout">
        <n-input-number 
          v-model:value="formData.rest.timeout" 
          :min="1"
          :max="60"
          placeholder="5"
        />
      </n-form-item>

      <n-divider title-placement="left">存档配置</n-divider>

      <n-form-item label="存档路径" path="save.path">
        <n-input 
          v-model:value="formData.save.path" 
          placeholder="/path/to/server/Pal/Saved"
        />
      </n-form-item>

      <n-form-item label="解码工具路径" path="save.decode_path">
        <n-input 
          v-model:value="formData.save.decode_path" 
          placeholder="解码工具路径（可选）"
        />
      </n-form-item>

      <n-form-item label="同步间隔(秒)" path="save.sync_interval">
        <n-input-number 
          v-model:value="formData.save.sync_interval" 
          :min="60"
          :max="3600"
          placeholder="120"
        />
      </n-form-item>

      <n-form-item label="备份间隔(秒)" path="save.backup_interval">
        <n-input-number 
          v-model:value="formData.save.backup_interval" 
          :min="3600"
          :max="86400"
          placeholder="14400"
        />
      </n-form-item>

      <n-form-item label="备份保留天数" path="save.backup_keep_days">
        <n-input-number 
          v-model:value="formData.save.backup_keep_days" 
          :min="1"
          :max="365"
          placeholder="7"
        />
      </n-form-item>

      <n-space justify="end" style="margin-top: 24px">
        <n-button @click="handleCancel">取消</n-button>
        <n-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '更新' : '创建' }}
        </n-button>
      </n-space>
    </n-form>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { NForm, NFormItem, NInput, NInputNumber, NButton, NSpace, NSwitch, NDivider } from 'naive-ui';
import { useMessage } from 'naive-ui';

const props = defineProps({
  server: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['submit', 'cancel']);

const message = useMessage();
const formRef = ref();
const submitting = ref(false);

const isEdit = computed(() => !!props.server?.id);

// 表单数据
const formData = ref({
  id: '',
  name: '',
  description: '',
  enabled: true,
  rcon: {
    address: '',
    password: '',
    use_base64: false,
    timeout: 5
  },
  rest: {
    address: '',
    username: 'admin',
    password: '',
    timeout: 5
  },
  save: {
    path: '',
    decode_path: '',
    sync_interval: 120,
    backup_interval: 14400,
    backup_keep_days: 7
  }
});

// 表单验证规则
const rules = {
  id: [
    { required: true, message: '请输入服务器ID', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_-]+$/, message: '服务器ID只能包含字母、数字、下划线和连字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入服务器名称', trigger: 'blur' }
  ],
  'rcon.address': [
    { required: true, message: '请输入RCON地址', trigger: 'blur' }
  ],
  'rcon.password': [
    { required: true, message: '请输入RCON密码', trigger: 'blur' }
  ],
  'rest.address': [
    { required: true, message: '请输入REST API地址', trigger: 'blur' }
  ],
  'rest.username': [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  'rest.password': [
    { required: true, message: '请输入REST API密码', trigger: 'blur' }
  ],
  'save.path': [
    { required: true, message: '请输入存档路径', trigger: 'blur' }
  ]
};

// 监听prop变化，更新表单数据
watch(() => props.server, (newServer) => {
  if (newServer) {
    formData.value = {
      id: newServer.id || '',
      name: newServer.name || '',
      description: newServer.description || '',
      enabled: newServer.enabled !== undefined ? newServer.enabled : true,
      rcon: {
        address: newServer.rcon?.address || '',
        password: newServer.rcon?.password || '',
        use_base64: newServer.rcon?.use_base64 || false,
        timeout: newServer.rcon?.timeout || 5
      },
      rest: {
        address: newServer.rest?.address || '',
        username: newServer.rest?.username || 'admin',
        password: newServer.rest?.password || '',
        timeout: newServer.rest?.timeout || 5
      },
      save: {
        path: newServer.save?.path || '',
        decode_path: newServer.save?.decode_path || '',
        sync_interval: newServer.save?.sync_interval || 120,
        backup_interval: newServer.save?.backup_interval || 14400,
        backup_keep_days: newServer.save?.backup_keep_days || 7
      }
    };
  } else {
    // 重置表单
    formData.value = {
      id: '',
      name: '',
      description: '',
      enabled: true,
      rcon: {
        address: '',
        password: '',
        use_base64: false,
        timeout: 5
      },
      rest: {
        address: '',
        username: 'admin',
        password: '',
        timeout: 5
      },
      save: {
        path: '',
        decode_path: '',
        sync_interval: 120,
        backup_interval: 14400,
        backup_keep_days: 7
      }
    };
  }
}, { immediate: true });

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value?.validate();
    submitting.value = true;
    emit('submit', formData.value);
  } catch (error) {
    message.error('请检查表单填写');
  } finally {
    submitting.value = false;
  }
};

// 取消操作
const handleCancel = () => {
  emit('cancel');
};
</script>

<style scoped>
.server-form {
  padding: 16px;
}
</style> 