<template>
  <div>
    <el-scrollbar
      id="logContainer"
      wrap-class="list"
      view-style="font-weight: bold;"
      view-class="view-box"
      :native="false"
    >
      <div>
        <ul class="infinite-list" id="lines-ui">
          <li v-for="line in lines" class="infinite-list-item">
            <Prism language="java" :code="line"></Prism>
          </li>
        </ul>
      </div>
    </el-scrollbar>
  </div>
</template>

<script>
import Prism from "vue-prismjs";
import "prismjs/themes/prism.css";
import { clearInterval, setInterval } from "timers";
export default {
  name: "tail-log",
  components: {
    Prism
  },
  props: {
    catalog: null,
    file: null
  },
  data() {
    return {
      lines: [],
      lineBuffer: [],
      tailWebSock: null,
      onSliding: true,
      roll: "un roll"
    };
  },
  created() {
    this.openTailLog();
    this.refreshLines();
  },
  beforeDestroy() {
    this.webSocketDest();
  },
  methods: {
    successNotify(msg) {
      const h = this.$createElement;
      this.$notify({
        title: "Success",
        message: h("i", { style: "color: #67C23A" }, msg),
        type: "success"
      });
    },
    warnNotify(msg) {
      const h = this.$createElement;
      this.$notify({
        title: "Warning",
        message: h("i", { style: "color: #E6A23C" }, msg),
        type: "warning"
      });
    },
    openTailLog() {
      if (
        this.catalog != null &&
        this.file != null &&
        this.catalog !== "" &&
        this.file !== ""
      ) {
        this.webSocketDest();
        const wsuri =
          "ws://" +
          window.location.host +
          "/api/tail/" +
          this.catalog +
          "/" +
          this.file;
        this.tailWebSock = new WebSocket(wsuri);
        this.tailWebSock.onmessage = this.webSocketOnMessage;
        this.tailWebSock.onopen = this.webSocketOnOpen;
        this.tailWebSock.onerror = this.webScoketOnError;
        this.tailWebSock.onclose = this.webScoketClose;
      }
    },
    webSocketOnOpen() {
      var lineUi = document.querySelector("#lines-ui");
      let action = { type: 0, ui_width: lineUi.clientWidth };
      this.webSocketSendJson(action);
      this.successNotify("Tail " + this.catalog + ":" + this.file + " success");
    },
    webScoketOnError() {
      this.warnNotify("Tail " + this.catalog + ":" + this.file + " failed");
    },
    webSocketOnMessage(e) {
      const redata = JSON.parse(e.data);
      switch (redata.type) {
        case 0: {
          break;
        }
        case 1: {
          // pile up a lot of data
          this.lineBuffer.push(redata.msg);
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
        let action = { type: 2 };
        this.webSocketSendJson(action);
      }, interval * 1000);
    },
    refreshLines() {
      clearInterval(this._refresh);
      this._refresh = setInterval(() => {
        var newLines = this.lineBuffer.splice(0, 200);
        if (newLines.length > 0) {
          var tempLines = this.lines.concat();
          tempLines.push(...newLines);
          if (tempLines.length > 500) {
            tempLines.splice(0, tempLines.length - 500);
          }
          this.lines = tempLines;
          this.sliding();
        }
      }, 1000); // ms
    },
    webSocketDest() {
      clearInterval(this._inter);
      clearInterval(this._refresh);
      if (this.tailWebSock != null) {
        this.tailWebSock.close();
        this.lines = [];
        this.tailWebSock = null;
      }
    },
    sliding() {
      if (this.onSliding) {
        this.$nextTick(() => {
          var container = this.$el.querySelector("#logContainer");
          container.scrollTop = container.scrollHeight;
        });
      }
    }
  }
};
</script>
<style>
.infinite-list-item {
  margin-top: 0px;
}

#logContainer {
  max-height: 680px;
}

.infinite-list {
  margin-bottom: 0px;
  padding-left: 20px;
}

code[class*="language-"],
pre[class*="language-"] {
  margin-top: 0px;
  margin-bottom: 0px;
  padding-top: 0px;
  padding-bottom: 0px;
  padding-left: 0px;
  padding-right: 0px;
}

li {
  list-style-type: none;
}

.el-scrollbar {
  overflow: scroll;
  max-height: 680px;
}
</style>
