<template>
  <div>
    <el-scrollbar
      wrap-class="list"
      view-style="font-weight: bold;"
      view-class="view-box"
      :native="false"
    >
      <div>
        <ul class="infinite-list">
          <li v-for="line in lines" class="infinite-list-item">{{ line }}</li>
        </ul>
      </div>
    </el-scrollbar>
  </div>
</template>

<script>
  import {clearInterval, setInterval} from "timers";

  export default {
    name: "tail-log",
    props: {
      catalog: null,
      file: null
    },
    data() {
      return {
        lines: [],
        tailWebSock: null
      };
    },
    created() {
      this.openTailLog();
    },
    beforeDestroy() {
      this.webSocketDest();
    },
    methods: {
      successNotify(msg) {
        const h = this.$createElement;
        this.$notify({
          title: "Success",
          message: h("i", {style: "color: #67C23A"}, msg)
        });
      },
      warnNotify(msg) {
        const h = this.$createElement;
        this.$notify({
          title: "Wran",
          message: h("i", {style: "color: #E6A23C"}, msg)
        });
      },
      openTailLog() {
        if (this.catalog != null && this.file != null && this.catalog != "" && this.file != "") {
          this.webSocketDest();
          const wsuri =
            "ws://127.0.0.1:3000/api/tail/" + this.catalog + "/" + this.file;
          this.tailWebSock = new WebSocket(wsuri);
          this.tailWebSock.onmessage = this.webSocketOnMessage;
          this.tailWebSock.onopen = this.webSocketOnOpen;
          this.tailWebSock.onerror = this.webScoketOnError;
          this.tailWebSock.onclose = this.webScoketClose;
        }
      },
      webSocketOnOpen() {
        let action = {type: 0};
        this.webSocketSendJson(action);
        this.successNotify("Tail " + this.catalog + " " + this.file + " success");
      },
      webScoketOnError() {
        this.warnNotify("Tail " + this.catalog + " " + this.file + " failed");
      },
      webSocketOnMessage(e) {
        const redata = JSON.parse(e.data);
        switch (redata.type) {
          case 0: {
            break;
          }
          case 1: {
            // pile up a lot of data
            this.lines.push(redata.msg);
            break;
          }
          case 2: {
            const heartInterval = parseInt(redata.msg);
            if (heartInterval > 0) {
              this.heartBeatSend(heartInterval);
            }
            break;
          }
          case 3: {
            break;
          }
        }
      },
      webSocketSendJson(data) {
        this.tailWebSock.send(JSON.stringify(data));
      },
      webScoketClose(e) {
        this.successNotify("Close " + this.catalog + " " + this.file);
        console.log("websocket close", e);
      },
      heartBeatSend(interval) {
        clearInterval(this._inter);
        this._inter = setInterval(() => {
          let action = {type: 2};
          this.webSocketSendJson(action);
        }, interval * 1000);
      },
      webSocketDest() {
        clearInterval(this._inter);
        if (this.tailWebSock != null) {
          this.tailWebSock.close();
          this.lines = [];
          this.tailWebSock = null;
        }
      }
    }
  };
</script>
<style>
  .infinite-list-item {
    margin-top: 5px;
  }
</style>
