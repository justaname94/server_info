import Vue from "vue";
import VueResource from "vue-resource";
import BootstrapVue from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import App from "./App.vue";

Vue.use(VueResource);
Vue.use(BootstrapVue);

Vue.http.options.root = "http://localhost:8000/";

new Vue({
  el: "#app",
  render: h => h(App)
});
