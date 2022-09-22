<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column';
import Badge from 'primevue/badge';
</script>

<script>
export default {
  data() {
    return {
      socket: io(
        // replace in dev mode with http://localhost:8080
        `${window.location.protocol}//${window.location.hostname}:${window.location.port}`,
        {
          withCredentials: false,
          transports: ['websocket'],
        }),
      messages: [],
      expanded: [],
    }
  },
  mounted() { },
  created() {
    this.socket.on('connect', () => {
      this.socket.emit('message_stack_req')
    })

    this.socket.on('message_stack_res', (data) => {
      this.messages = JSON.parse(data)
    })

    this.socket.on('message_update', (data) => {
      this.messages = [JSON.parse(data), ...this.messages]
    })
  },
  methods: {
    ellipis: (target, limit) => {
      return target.length > limit ? `${target.substring(0, limit - 3)}...` : target
    }
  }
}
</script>

<template>
  <DataTable :value="messages" :paginator="true" :rows="20" :row-hover="true" v-model:expandedRows="expanded"
    data-key="uid">
    <Column :expander="true" headerStyle="width: 3rem"> </Column>
    <Column field="updated" header="Updated" :sortable="true">
      <template #body="slotProps">
        {{new Date(slotProps.data.updated).toLocaleString()}}
      </template>
    </Column>
    <Column field="otp" header="OTP">
      <template #body="slotProps">
        <Badge v-if="slotProps.data.otp" severity="success" size="large">
          {{slotProps.data.otp}}
        </Badge>
        <Badge v-else severity="warning" size="large" value="Not found"></Badge>
      </template>
    </Column>
    <Column field="sender" header="Sender"></Column>
    <Column field="recipient" header="Recipient"></Column>
    <Column field="message" header="Message">
      <template #body="slotProps">
        {{ellipis(slotProps.data.message, 50)}}
      </template>
    </Column>
    <template #expansion="slotProps">
      Original message
      <div>
        <pre>
          {{slotProps.data}}
        </pre>
      </div>
    </template>
  </DataTable>
</template>

<style scoped>
pre {
  white-space: pre-line
}
</style>
