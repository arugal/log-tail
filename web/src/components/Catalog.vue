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
      id="dialog"
      :title='currentCatalog+" - "+currentFile'
      :visible.sync="dialogTableVisible"
      @closed="closeTailLog"
    >
      <tail-log :catalog="currentCatalog" :file="currentFile" :key="tailLogKey"></tail-log>
    </el-dialog>
  </div>
</template>

<script>
import tailLog from "./Taillog.vue";

export default {
  data() {
    return {
      catalogs: [],
      dialogTableVisible: false,
      lines: [],
      currentCatalog: null,
      currentFile: null,
      tailLogKey: 1
    };
  },
  created() {
    this.fetchData();
  },
  components: {
    tailLog
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
    fetchData() {
      fetch("/api/catalog", {
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
      this.currentCatalog = catalog;
      this.currentFile = file;
      ++this.tailLogKey;
      this.dialogTableVisible = true;
    },
    closeTailLog() {
      this.currentCatalog = "";
      this.currentFile = "";
      ++this.tailLogKey;
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
.el-dialog__body {
  padding-top: 0px;
  padding-left: 0px;
  padding-bottom: 0px;
}
</style>
