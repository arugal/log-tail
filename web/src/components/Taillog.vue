<template>
  <div class="test"></div>
</template>

<script>
import { setInterval, clearInterval } from "timers";
export default {
  data() {
    return {
      websock: null
    };
  },
  created() {
    this.initWebSocket();
    this.heartbeatInterval();
  },
  destroyed() {
    this.websock.close(); //离开路由之后断开websocket连接
    clearInterval(this._inter);
  },
  methods: {
    initWebSocket() {
      //初始化weosocket
      const wsuri = "ws://127.0.0.1:3000/api/tail/root/3.1.json";
      this.websock = new WebSocket(wsuri);
      this.websock.onmessage = this.websocketonmessage;
      this.websock.onopen = this.websocketonopen;
      this.websock.onerror = this.websocketonerror;
      this.websock.onclose = this.websocketclose;
    },
    websocketonopen() {
      // 连接建立之后执行send方法发送数据
      let actions = { type: 0 };
      this.websocketsend(JSON.stringify(actions));
    },
    websocketonerror() {
      //连接建立失败重连
      this.initWebSocket();
    },
    websocketonmessage(e) {
      //数据接收
      const redata = JSON.parse(e.data);
      console.log(redata);
    },
    websocketsend(Data) {
      //数据发送
      this.websock.send(Data);
    },
    websocketclose(e) {
      //关闭
      console.log("断开连接", e);
    },
    heartbeatInterval() {
      this._inter = setInterval(() => {
        let actions = { type: 2 };
        this.websocketsend(JSON.stringify(actions));
      }, 5000);
    }
  }
};
</script>
<style>
</style>