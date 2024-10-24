import { createRouter, createWebHistory } from 'vue-router';
import BillingDashboard from '../components/BillingDashboard.vue'; // Update the path
import ViewBills from '../views/ViewBills.vue';
import AddBill from '../views/AddBill.vue';
import ManageBills from '../views/ManageBills.vue';

const routes = [
  { path: '/', component: BillingDashboard }, // BillingDashboard as the default route
  { path: '/view-bills', component: ViewBills },
  { path: '/add-bill', component: AddBill },
  { path: '/manage-bills', component: ManageBills },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
