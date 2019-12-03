<template>
  <div class="cards">
    <b-card
      header="xataka.com"
      title="Xataka - Tecnología y gadgets, móviles, informática, electrónica"
    >
      <b-list-group>
        <b-list-group-item>
          <strong>Grade:</strong>
          <b-badge>A</b-badge>
        </b-list-group-item>
        <b-list-group-item>
          <strong>Previous grade:</strong>
          <b-badge>B</b-badge>
        </b-list-group-item>
        <b-list-group-item>
          <strong>Logo:</strong>
          <img height="36px" width="36px" src="../assets/logo.png" alt />
        </b-list-group-item>
        <b-list-group-item>Is not down</b-list-group-item>
        <b-list-group-item>Have not changed lately</b-list-group-item>
      </b-list-group>
      <hr />
      <h4>Servers</h4>
      <div>
        <b-table striped :items="items"></b-table>
      </div>
    </b-card>
  </div>
</template>

<script>
export default {
  data() {
    return {
      items: [
        {
          address: "13.35.125.83",
          grade: "A",
          country: "US",
          owner: "Amazon Technologies Inc. (AT-88-Z)"
        },
        {
          address: "13.35.125.87",
          grade: "A",
          country: "US",
          owner: "Amazon Technologies Inc. (AT-88-Z)"
        },
        {
          address: "13.35.125.4",
          grade: "A",
          country: "US",
          owner: "Amazon Technologies Inc. (AT-88-Z)"
        },
        {
          address: "13.35.125.72",
          grade: "A",
          country: "US",
          owner: "Amazon Technologies Inc. (AT-88-Z)"
        }
      ]
    };
  }
};
</script>

<style scoped>
</style>