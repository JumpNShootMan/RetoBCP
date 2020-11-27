<template>
  <v-row>
    <v-col class="text-center">
      <img
        src="../assets/bcp1.jpg"
        alt="Vuetify.js"
        class="mb-5" 
      >
      <blockquote class="blockquote">
        &#8220;Bienvenidos a la página principal del Reto BCP!.&#8221;
        <footer>
          <small>
            <em>&mdash;Por Daniel Núñez y Christian Espíritu</em>
          </small>
        </footer>
      </blockquote>
      <!-- <p v-for="trans in transactions" :key="trans.ID">
        
        {{trans.ID}}
        {{trans.CreatedAt}}
        {{trans.UpdatedAt}}
        {{trans.DeletedAt}}
        {{trans.from_user}}
        {{trans.amount}}
        {{trans.description}}
      
      </p> -->
        <v-simple-table dark>
    <template v-slot:default>
      <thead>
        <tr>
          <th class="text-center">
            ID
          </th>
          <th class="text-center">
            Monto
          </th>
          <th class="text-center">
            Fecha Transacción
          </th>
          <th class="text-center">
            Usuario Origen
          </th>
          <th class="text-center">
            Usuario Destino
          </th>
          <th class="text-center">
            Categoría
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="item in transactions"
          :key="item.ID"
        >
          <td>{{ item.ID }}</td>
          <td>{{ item.amount }}</td>
          <td>{{ item.CreatedAt}}</td>
          <td>{{ item.From.ID}}</td>
          <td>{{ item.To.ID}}</td>
          <td>{{ item.category_id}}</td>

          

        </tr>
      </tbody>
    </template>
  </v-simple-table>
    </v-col>
  </v-row>
</template>
<script lang="ts">
import Vue from 'vue'
import axios from 'axios'

type Transaction = {
  ID: number,
  CreatedAt: Date,
  UpdatedAt: Date,
  DeletedAt: Date,
  amount: number,
  from_user: number,
  description: number;
  category_id: number;
}
export default Vue.extend({
  async mounted(){ //en startup; async espera a que termine una instrucción para ir a la siguiente :) ☺
    const response = await axios.get<Transaction[]>('http://localhost:1996/api/v1/transactions') //[] es nueva forma de declarar arrays
    console.log(response.data)
    this.transactions = response.data
  },
  data(){
    return{
      transactions: Array<Transaction>()
      
    };
  }
})
</script>
