<template>
  <b-card :header="website.domain" :title="website.title" class="mt-2">
    <b-list-group>
      <b-list-group-item>
        <strong>Grade:</strong>
        <b-badge>{{ website.grade }}</b-badge>
      </b-list-group-item>
      <b-list-group-item v-if="website.previousGrade">
        <strong>Previous grade:</strong>
        <b-badge>{{ website.previousGrade }}</b-badge>
      </b-list-group-item>
      <b-list-group-item>
        <strong>Logo:</strong>
        <img height="36px" width="36px" v-if="!website.isDown" :src="logo" alt />
      </b-list-group-item>
      <b-list-group-item>
        <span v-if="website.isDown">Is Down</span>
        <span v-else>Is not down</span>
      </b-list-group-item>
      <b-list-group-item>Have not changed lately</b-list-group-item>
    </b-list-group>
    <hr />
    <h4>Servers</h4>
    <div>
      <b-table striped :items="website.servers"></b-table>
    </div>
  </b-card>
</template>

<script>
export default {
  props: {
    website: {
      type: Object,
      required: true
    }
  },
  computed: {
    logo: function() {
      return `${this.$http.options.root.slice(0, -1)}${this.website.logo}`;
    }
  }
};
</script>