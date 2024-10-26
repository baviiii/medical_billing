import { createRouter, createWebHistory } from 'vue-router';
import BillingDashboard from '../components/BillingDashboard.vue'; // Update the path
import ViewBills from '../views/ViewBills.vue';
import AddBill from '../views/AddBill.vue';
import ManageBills from '../views/ManageBills.vue';
import UserLogin from '../views/UserLogin.vue';
import UserRegister from '../views/UserRegister.vue';

const routes = [
  { path: '/', component: BillingDashboard }, // BillingDashboard as the default route
  { path: '/view-bills', component: ViewBills },
  { path: '/add-bill', component: AddBill },
  { path: '/manage-bills', component: ManageBills },
  { path: '/login', component: UserLogin },
  { path: '/register', component: UserRegister },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Route guard
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token'); // Check for token
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login'); // Redirect to login if not authenticated
  } else {
    next(); // Proceed to the route
  }
});

export default router;

