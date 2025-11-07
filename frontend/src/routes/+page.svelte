<script>
  import { goto } from '$app/navigation';
  import Header from '$lib/components/layout/Header.svelte';
  import Footer from '$lib/components/layout/Footer.svelte';
  import camera from '$lib/assets/camera.jpg';
  import roteador from '$lib/assets/roteador.png';
  import alarme from '$lib/assets/alarme.jpg';
  import background from '$lib/assets/placasolar2.jpg';

  const goToChat = () => goto('/chat');
</script>

<Header />

<main class="pt-28 bg-white text-gray-700">

<section class="relative min-h-screen flex items-center justify-center bg-cover bg-center -mt-7" style="background-image: url('{background}')">
  
  <div class="absolute inset-0 bg-black opacity-60 z-10"></div>

  <div class="relative z-20 max-w-2xl text-center px-4">
    <h1 class="text-4xl sm:text-5xl font-bold text-white mb-4">
      Conheça a Iva
    </h1>
    <p class="text-lg leading-relaxed text-white mb-10">
      A <strong>Iva</strong> (Intelbras Virtual Assistant)  é sua assistente de inteligência artificial desenvolvida para lhe ajudar a solucionar dúvidas sobre produtos e soluções <strong>Intelbras</strong> — de forma rápida e confiável.
    </p>
    <button
      on:click={goToChat}
      class="bg-[#00995D] hover:bg-[#00814F] text-white px-10 py-4 rounded-xl font-semibold text-lg transition-colors duration-300 cursor-pointer"
    >
      Falar com a Iva
    </button>
  </div>
</section>

  <section class="bg-gray-50 py-16 px-6">
    <div class="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-10 items-center">
      <div>
        <h2 class="text-3xl font-bold text-[#00995D] mb-4">Suporte técnico em tempo real</h2>
        <p class="text-lg leading-relaxed">
         Precisa configurar uma câmera, ajustar um roteador ou entender a instalação de um sistema de alarme?
          A Intelbras Virtual Assistant fornece <strong>orientações precisas</strong> e 
          <strong>passo a passo detalhado</strong> em segundos.
        </p>
      </div>
      <div class="flex justify-center">
        <img src={camera} alt="Câmera Intelbras" class="rounded-xl shadow-md max-w-full h-auto" />
      </div>
    </div>
  </section>

  <section class="py-16 px-6">
    <div class="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-10 items-center">
      <div class="order-2 md:order-1 flex justify-center">
        <img src={roteador} alt="Roteador Intelbras" class="rounded-xl shadow-md max-w-full h-auto" />
      </div>
      <div class="order-1 md:order-2">
        <h2 class="text-3xl font-bold text-[#00995D] mb-4">Aprendizado contínuo</h2>
        <p class="text-lg leading-relaxed">
          Nossa assistente evolui com base nas suas interações. Quanto mais você usa a Iva, mais ela entende suas necessidades e oferece respostas personalizadas sobre os produtos Intelbras.
        </p>
      </div>
    </div>
  </section>

  <section class="bg-gray-50 py-16 px-6">
    <div class="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-10 items-center">
      <div>
        <h2 class="text-3xl font-bold text-[#00995D] mb-4">Feito para quem está no campo</h2>
        <p class="text-lg leading-relaxed">
          A Intelbras Virtual Assistant foi pensada também para <strong>técnicos, instaladores e integradores</strong> que precisam de suporte rápido, direto do celular ou 
          notebook — sem perder tempo em filas de atendimento.
        </p>
      </div>
      <div class="flex justify-center">
        <img src={alarme} alt="Alarme Intelbras" class="rounded-xl shadow-md max-w-full h-auto" />
      </div>
    </div>
  </section>

  <Footer />

</main>
