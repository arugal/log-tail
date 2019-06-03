<template>
  <div id="app">
    <header class="grid-content header-color">
      <el-row>
        <svg>
          <use xlink:href="#icon-tail"></use>
        </svg>
        <a class="brand" href="#">tail</a>
      </el-row>
    </header>
    <section>
      <el-row :gutter="20">
        <el-col id="side-nav" :xs="24" :md="4">
          <el-menu default-active="1" mode="vertical" theme="light" router @select="handleSelect">
            <el-menu-item index="/">Catalog</el-menu-item>
            <el-menu-item index>Help</el-menu-item>
          </el-menu>
        </el-col>

        <el-col :xs="24" :md="20">
          <div id="content">
            <router-view></router-view>
          </div>
        </el-col>
      </el-row>
    </section>
    <footer></footer>
  </div>
</template>

<script>
export default {
  methods: {
    handleSelect(key, path) {
      if (key === "") {
        window.open("https://github.com/Arugal/log-tail");
      }
    }
  }
};
</script>

<style>
body {
  background-color: #fafafa;
  margin: 0px;
  font-family: -apple-system, BlinkMacSystemFont, Helvetica Neue, sans-serif;
}

header {
  width: 100%;
  height: 60px;
}

.header-color {
  background: #606266;
}

#content {
  margin-top: 20px;
  padding-right: 40px;
}

.brand {
  color: #fff;
  background-color: transparent;
  margin-left: 20px;
  float: left;
  line-height: 25px;
  font-size: 25px;
  padding: 15px 15px;
  height: 30px;
  text-decoration: none;
}
</style>
