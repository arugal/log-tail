<template>
  <div>
    <div v-for="item in catalogs" :key="item.catalog">
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span>{{item.catalog}}</span>
        </div>
        <div>
          <el-button
            plain
            v-for="file in item.child_file"
            :key="file"
            @click="openTailLog(item.catalog, file)"
          >{{file}}</el-button>
        </div>
      </el-card>
    </div>
    <el-dialog
      :title='currentCatalog+" - "+currentFile'
      :visible.sync="dialogTableVisible"
      @closed="webSocketDest"
    >
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
    </el-dialog>
  </div>
</template>

<script>
import { setInterval, clearInterval } from "timers";
export default {
  data() {
    return {
      catalogs: [],
      tailWebSock: null,
      dialogTableVisible: false,
      lines: [],
      currentCatalog: null,
      currentFile: null
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    successNotify(msg) {
      const h = this.$createElement;
      this.$notify({
        title: "Success",
        message: h("i", { style: "color: #67C23A" }, msg)
      });
    },
    warnNotify(msg) {
      const h = this.$createElement;
      this.$notify({
        title: "Wran",
        message: h("i", { style: "color: #E6A23C" }, msg)
      });
    },
    tailLog() {},
    fetchData() {
      fetch("http://127.0.0.1:3000/api/catalog", {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"
        }
      })
        .then(res => {
          return res.json();
        })
        .then(json => {
          this.catalogs = json;
          this.successNotify("Get api catalog success");
        })
        .catch(err => {
          this.warnNotify("Get api catalog failed");
        });
    },
    openTailLog(catalog, file) {
      if (this.tailWebSock != null) {
        this.webSocketDest();
      }
      const wsuri = "ws://127.0.0.1:3000/api/tail/" + catalog + "/" + file;
      this.tailWebSock = new WebSocket(wsuri);
      this.tailWebSock.onmessage = this.webSocketOnMessage;
      this.tailWebSock.onopen = this.webSocketOnOpen;
      this.tailWebSock.onerror = this.webScoketOnError;
      this.tailWebSock.onclose = this.webScoketClose;
      this.heartBeatSend();
      this.currentCatalog = catalog;
      this.currentFile = file;
      this.dialogTableVisible = true;
    },
    webSocketOnOpen() {
      let action = { type: 0 };
      this.webSocketSendJson(action);
      this.successNotify(
        "Tail " + this.currentCatalog + " " + this.currentFile + " success"
      );
    },
    webScoketOnError() {
      this.warnNotify(
        "Tail " + this.currentCatalog + " " + this.currentFile + " failed"
      );
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
    webSocketSend(data) {
      this.tailWebSock.send(data);
    },
    webScoketClose(e) {
      this.successNotify(
        "Close " + this.currentCatalog + " " + this.currentFile
      );
      this.currentCatalog = "";
      this.currentFile = "";
      console.log("websocket close", e);
    },
    heartBeatSend(interval) {
      clearInterval(this._inter);
      this._inter = setInterval(() => {
        let action = { type: 2 };
        this.webSocketSendJson(action);
      }, interval * 1000);
    },
    webSocketDest() {
      clearInterval(this._inter);
      this.tailWebSock.close();
      this.lines = [];
    }
  }
};
</script>

<style>
.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both;
}

.box-card {
  width: 100%;
}

.el-button {
  margin-top: 10px;
}

.el-dialog {
  width: 90%;
  height: 80%;
}

.list {
  max-height: 600px;
}

.infinite-list-item {
  margin-top: 5px;
}
</style>
