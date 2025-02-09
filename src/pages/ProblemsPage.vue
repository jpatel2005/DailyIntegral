<template>
  <LoadingCogs v-if="isLoading" />
  <div v-else-if="problemLatex === null" class="w-full text-center">
    <div class="m-auto">
      <h3 class="text-white font-semibold text-5xl tracking-wide mt-8 lowercase">{{ dateStr }}</h3>
      <h3 class="text-white font-medium text-4xl md:text-8xl tracking-wide mt-24 md:mt-64">No problem available!!!</h3>
    </div>
  </div>
  <div v-else class="flex justify-center items-center w-full h-full overflow-x-hidden">
    <div class="flex flex-col items-center w-full h-full">
      <h3 class="text-white font-semibold text-3xl lg:text-5xl tracking-wide mt-4 lg:mt-8 lowercase">{{ dateStr }}</h3>
      <div class="w-full flex md:inline-block md:mt-[-2em] lg:mt-[-3em]">
        <div v-if="!isDaily" class="text-lg float-left w-32 text-center ml-4 px-4 py-3 w-fit">
          <button class="text-white bg-blue-600 hover:bg-blue-700 px-3 py-2 rounded-md" @click="open">
            Solution 💡
          </button>
        </div>
        <div :class="problemState || problemState === 2 ? 'visible' : 'invisible'"
          class="text-white text-lg bg-gray-600 border-2 border-gray-800 float-right w-32 text-center rounded-lg items-center justify-center flex flex-row gap-x-2 shadow-xl mr-10 px-4 py-3 w-fit max-sm:m-auto">
          <span>{{ problemState === 1 ? 'Solved' : problemState === 2 ? 'Solved (as daily problem)' : '\u200b' }}</span>
          <span class="pi pi-check-circle text-green-500"></span>
        </div>
      </div>

      <div class="w-full md:px-10 lg:px-20 md:mt-8 lg:mt-8 z-50 flex items-center justify-between">
        <div class="hidden sm:inline">
          <Stopwatch />
        </div>
        <div class="m-auto md:m-0">
          <div class="font-sans bg-gray-900 px-6 py-4 rounded-lg">
            <h2 class="text-center text-2xl font-bold mb-4 text-white uppercase">
              Important Notes:
            </h2>
            <div class="overflow-y-scroll max-w-96 md:max-w-[32em] max-h-28 md:max-h-32">
              <ul class="list-disc text-white ml-6 lg:text-base mdtext-sm">
                <li>For indefinite integrals, don't include the constant (+ c)</li>
                <li>Logarithms are base e unless specified otherwise</li>
                <li>
                  When multiplying variables (or expressions) together, use a
                  multiplication symbol&nbsp;"&bull;"&nbsp;(e.g. xy &rarr; x&nbsp;&bull;&nbsp;y)
                </li>
                <li>Always use parentheses with trigonometric functions (e.g. cos(x) instead of cosx), logarithms, and
                  etc</li>
                <li>Use fractional exponents [e.g. x^(3/2)] instead of roots or nested powers [e.g. sqrt(x^3) or
                  (x^3)^(1/2)]</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
      <div
        class="equation-container !mb-8 px-12 select-none overflow-x-scroll scrollbar-thin scrollbar-thumb-gray-500 scrollbar-track-gray-700 w-full">
        <vue-latex :expression="problemLatex" :fontsize="screenSize" class="text-white text-4xl md:text-6xl lg:text-8xl"
          display-mode />
      </div>
      <div class="text-center">
        <math-field id="math-input" ref="math-input"
          class="items-center inline-block w-[28em] md:w-[32em] lg:w-[48em] max-w-[90vw] md:text-2xl lg:text-4xl py-1 md:py-1.5 px-1.5 md:px-2 pr-0"
          defaultMode="math" :smartFence="true" :smartSuperscript="false" :smartMode="false"
          @input='event => answer = event.target.value'></math-field>
      </div>
      <div class="mt-8 items-center flex flex-row gap-x-3">
        <button
          class="text-white float-left focus:outline-none font-medium text-lg rounded-lg text-xl px-4 py-2 text-center bg-blue-600 disabled:bg-gray-300 hover:bg-blue-700 border-transparent"
          @click="checkAnswer" :disabled="checkAnswerActive">Check Answer</button>
        <span class="pi pi-spin pi-cog text-3xl inline text-white" v-show="checkAnswerActive"></span>
      </div>
      <ConfettiExplosion v-if="visible" />
      <div class="items-center flex flex-row gap-x-2 mt-3.5" :class="answerStatus ? 'text-green-400' : 'text-red-600'">
        <p class="text-2xl font-serif tracking-wide" v-show="answerStatus !== null">{{ answerMessage }}</p>
        <span class="pi text-lg mt-[-2px]" :class="answerStatus ? 'pi-check-circle' : 'pi-times-circle'"
          v-show="answerStatus !== null"></span>
      </div>
    </div>
    <div ref="keyboardContainer" class="keyboard-footer"></div>
  </div>
  <ModalsContainer />
</template>

<style scoped>
#math-input:focus-within {
  outline: 4px solid #1E48A4;
  border-radius: 4px;
}

.equation-container {
  margin-top: .5em;
}

.keyboard-footer {
  position: absolute;
  left: 0;
  width: 100%;
  height: fit-content;
}
</style>

<style>
.katex-display {
  margin-top: .5em;
  margin-bottom: .25em;
}
</style>

<script>
import { onMounted, ref, nextTick, useTemplateRef } from 'vue';
import { VueLatex } from 'vatex';
import { useRoute } from 'vue-router';
import { useAuth0 } from '@auth0/auth0-vue';
import ConfettiExplosion from "vue-confetti-explosion";
import 'primeicons/primeicons.css'
import { convertDateToUTC, formatDateISO, formatDateString, sleep } from '../time';
import router from '../router';
import { ModalsContainer, useModal } from 'vue-final-modal';
import SolutionModal from '../components/SolutionModal.vue';
import LoadingCogs from '../components/LoadingCogs.vue';
import Stopwatch from '../components/Stopwatch.vue';
export default {
  setup() {
    const { getAccessTokenSilently } = useAuth0();
    const isLoading = ref(true);
    const dateStr = ref('');
    const problemLatex = ref('');
    const solutionLatex = ref('');
    const stepsLatex = ref('');
    const answer = ref('');
    const visible = ref(false);
    const isDaily = ref(false);
    const userData = ref(null);
    const problemState = ref(0);
    const checkAnswerActive = ref(false);
    const formattedDate = ref('');
    const answerStatus = ref(null);
    const answerMessage = ref('');
    const mathInput = useTemplateRef('math-input');
    const keyboardContainer = ref(null);
    // render confetti explosion
    const explode = async function () {
      visible.value = false;
      await nextTick();
      visible.value = true;
    }
    // Configure date related information (use router param if any)
    const route = useRoute();
    const curr = ref('');

    // Perform the appropriate actions depending on whether or not the given answer was correct
    const manageAnswerStatus = async function (status, message) {
      // Answer status has been determined so set to false
      checkAnswerActive.value = false;
      answerStatus.value = status;
      // Use message if one is provided, or else set based on status
      answerMessage.value = message ?? (status ? 'Correct!' : 'Incorrect!');
      if (status) {
        // Confetti animation
        await explode();
        problemState.value = await problemSolveStatus(userData.value);
      }
    }
    // Retrieve the user's state for the current problem (unsolved = 0, solved not as daily problem = 1, solved as daily problem = 2)
    const problemSolveStatus = async function (user) {
      if (user?.username === undefined) {
        return 0;
      }
      const date = formattedDate.value.trim();
      try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/problemStatus/${date}?username=${user?.username ?? ''}`, {
          method: 'GET',
        });
        if (!response.ok) {
          return 0;
        }
        const data = await response.json();
        return [0, 1, 2].includes(data) ? data : 0;
      } catch {
        return 0;
      }
    }
    // Retrieve problem as latex given a date string (YYYY-MM-DD)
    const getProblem = async function () {
      const date = formattedDate.value;
      try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/problem/${date}`, {
          method: 'GET',
        });
        if (!response.ok) {
          return null;
        }
        const data = await response.json();
        if (data.problem === undefined) {
          return null;
        }
        return data.problem;
      } catch {
        return null;
      }
    }
    // Retrieve solution and steps (steps may just be an empty string) given a date string (YYYY-MM-DD)
    const getSolution = async function () {
      const date = formattedDate.value;
      try {
        const token = await getAccessTokenSilently();
        const response = await fetch(`${import.meta.env.VITE_API_URL}/solution/${date}`, {
          method: 'GET',
          headers: {
            'Authorization': token,
          }
        });
        if (!response.ok) {
          return [null, null];
        }
        const json = await response.json();
        return [json.answer, json.steps || ''];
      } catch {
        // Return null array to preserve destructuring
        return [null, null];
      }
    }
    // Send request to the api to verify if the answer is correct (if so, then respond accordingly)
    const checkAnswer = async function () {
      checkAnswerActive.value = true;
      answerStatus.value = null;
      // sleep so that the user has to wait a bit longer
      await sleep(2000);
      // parse the input (replace cdot multiplication with asterik for proper parsing)
      const input = answer.value.trim()
        .replace(/\\cdot/g, '*');
      if (input === '') {
        manageAnswerStatus(false);
        return false;
      }
      // First check if answer is invalid (e.g. contains integral)
      if (input.includes('\\int')) {
        manageAnswerStatus(false, 'Answer cannot contain an integral!');
        return false;
      }
      // Check if answer contains comparison operators
      if (input.includes('>') || input.includes('<') || input.includes('=')) {
        manageAnswerStatus(false, 'Answer cannot contain comparison operators!');
        return false;
      }
      // First verify that the user is signed in (no verification otherwise)
      if (userData.value?.username === undefined) {
        manageAnswerStatus(false, 'Please sign-in before submitting!');
        return false;
      }
      const date = formattedDate.value;
      // Check if the answer is correct using verify endpoint
      try {
        const urlEncoded = new URLSearchParams();
        urlEncoded.append('answer', input);
        urlEncoded.append('username', userData.value.username);
        const token = await getAccessTokenSilently();
        const response = await fetch(`${import.meta.env.VITE_API_URL}/verify/${date}`, {
          method: 'POST',
          headers: {
            'Authorization': token,
            'Content-Type': 'application/x-www-form-urlencoded',
          },
          body: urlEncoded.toString()
        });
        // Check if request is being rate limited
        console.log(response.status);
        if (response.status === 429) {
          manageAnswerStatus(false, 'Please wait a moment before submitting!')
          return false;
        }
        if (!response.ok) {
          manageAnswerStatus(false);
          return false;
        }
        const data = await response.json();
        const status = data === 'true';
        manageAnswerStatus(status);
        return status;
      } catch {
        manageAnswerStatus(false);
        return false;
      }
    }
    // Setup solution modal
    const { open, close } = useModal({
      component: SolutionModal,
      attrs: {
        problem: problemLatex,
        solution: solutionLatex,
        steps: stepsLatex,
        onConfirm() {
          close()
        },
      },
    });
    onMounted(async () => {
      const { user } = useAuth0();
      const date = new Date();
      userData.value = user.value;
      isLoading.value = true;
      isDaily.value = false;
      // Use router param for date, or else use current date (active daily problem)
      curr.value = route.params['pdate'] !== '' ? convertDateToUTC(new Date(route.params['pdate'])) : date;
      // If date is in the future, then route to home page (prevent trolling)
      if (route.params['pdate'] && (!/^\d{4}-\d{2}-\d{2}$/.test(route.params['pdate']) || convertDateToUTC(new Date()).getTime() - curr.value.getTime() < 0)) {
        await sleep(250);
        router.push('/');
        return;
      }
      formattedDate.value = formatDateISO(curr.value);
      dateStr.value = formatDateString(curr.value);
      // sleep so that the page will load
      await sleep(250);
      // Get problem and render the equation
      // If problem latex is null, then render the page differently (say "no problem available")
      problemLatex.value = await getProblem();
      if (problemLatex.value === null) {
        isLoading.value = false;
        return;
      }
      // Get the user's existing status of solving the problem
      problemState.value = await problemSolveStatus(userData.value);
      // show 'solution' button only if problem is not the daily problem (check if routing indicates that this problem is a daily)
      // either no route param (daily) or the route param is the daily's date
      if (!route.params['pdate'] || formatDateISO(date) === formatDateISO(curr.value)) {
        isDaily.value = true;
      } else {
        [solutionLatex.value, stepsLatex.value] = await getSolution();
      }
      isLoading.value = false;
      // sleep for a bit so that the math input can mount
      await sleep(1000);
      // Remove undesired menu items from math-input (do this after loading to after errors -- nothing should change visually on render so this should be fine, plus the operation will be quick)
      const itemsToRemove = new Set(['paste', 'ce-evaluate', 'ce-simplify', 'ce-solve', 'mode', 'variant', 'color', 'background-color', 'insert-matrix']);
      mathInput.value.menuItems = mathInput.value.menuItems.filter(item => !itemsToRemove.has(item.id));
    });
    // Manage resizing of the latex expression
    const screenSize = ref('default');
    const updateScreenSize = function () {
      const width = window.innerWidth;
      if (width >= 1024) {
        screenSize.value = 72;
      } else if (width >= 768) {
        screenSize.value = 40;
      } else {
        screenSize.value = 18;
      }
      screenSize.value *= (7 / 8);
    };
    const isMobile = function () {
      let check = false;
      (function (a) { if (/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows ce|xda|xiino/i.test(a) || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(a.substr(0, 4))) check = true; })(navigator.userAgent || navigator.vendor || window.opera);
      return check;
    }
    onMounted(() => {
      // manage dynamic sizing of the rendered latex
      updateScreenSize();
      window.addEventListener('resize', updateScreenSize);
    });
    onMounted(async () => {
      // Sleep while the page is loading (retrieving the problem + rendering page)
      while (isLoading.value) {
        await sleep(500);
      }
      // sleep for a bit extra then set the math keyboard container
      await sleep(250);
      if (keyboardContainer.value && window.mathVirtualKeyboard) {
        window.mathVirtualKeyboard.container = keyboardContainer.value;
        keyboardContainer.value.style.bottom = '0';
        const renderKeyboard = function () {
          if (!keyboardContainer.value) {
            return;
          }
          const height = Number(window.mathVirtualKeyboard.boundingRect.height);
          // configure the keyboard container (should work decently for mobile)
          keyboardContainer.value.style.height = height + 'px';
          keyboardContainer.value.style.marginBottom = '-' + (window.innerWidth >= 1440 ? .5 * height : (isMobile() && window.innerWidth >= window.innerHeight ? Math.max(2.5, window.innerWidth / window.innerHeight) * height : (isMobile() ? height * Math.max(1.5, Math.min(1.5, window.innerHeight / window.innerWidth)) : height))) + 'px';
          document.documentElement.scrollTop = document.documentElement.scrollHeight;
        }
        // create an event listener that will detect when the keyboard is opened (hopefully) and update a ref storing the keyboard height
        window.mathVirtualKeyboard.addEventListener("geometrychange", renderKeyboard);
        // also have this occur when the window is resized
        window.mathVirtualKeyboard.addEventListener("resize", renderKeyboard);
      }
    });
    return { dateStr, problemLatex, answer, answerMessage, answerStatus, checkAnswer, checkAnswerActive, visible, mathInput, problemState, isDaily, isLoading, open, screenSize, keyboardContainer };
  },
  components: {
    VueLatex, ConfettiExplosion, ModalsContainer, LoadingCogs, Stopwatch
  },
}
</script>
