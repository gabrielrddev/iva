<script>
import { goto } from '$app/navigation';
import logo from '$lib/assets/logo.png';
const goToAccount = () => goto('/account');
const goToHome = () => goto('/');
</script>
<div>
<header class="w-full flex items-center justify-between px-6 py-4 bg-white shadow-sm fixed top-0 left-0 z-50">
  <div class="flex items-center gap-3 cursor-pointer">
    <button class="cursor-pointer" on:click={goToHome}>
      <img src={logo} alt="Intelbras" class="w-32" />
    </button>
  </div>

  <button
    on:click={goToAccount}
    class="bg-[#00995D] hover:bg-[#00814F] text-white w-11 h-11 flex items-center justify-center rounded-full transition-colors duration-300 cursor-pointer"
    aria-label="Minha Conta"
  >
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-6 h-6">
      <path stroke-linecap="round" stroke-linejoin="round"
        d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.5 20.25a8.25 8.25 0 0115 0" />
    </svg>
  </button>
</header>
</div>