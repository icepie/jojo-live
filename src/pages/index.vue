<script setup lang="ts">
import NPlayer from "nplayer";
import Hls from "hls.js";
// import Flv from "flv.js";
// import NPlayer from "@nplayer/vue/";

import axios from "axios";

import { ABtn, ADialog, ACard } from "anu-vue";

import { useToast } from "vue-toastification";

import { Waline } from "@waline/client/component";

import "@waline/client/dist/waline.css";
import { computed } from "vue";
import { useRoute } from "vue-router";

const showDialog = ref(false);

const showSleepDialog = ref(false);

const serverURL = "https://icepie.singzer.cn";

const hlsUrl = "https://ice.singzer.cn/live/jojo.m3u8";

const player = new NPlayer({
  poster:
    "https://camo.githubusercontent.com/9dc29b67b014909713c12a3cf0f5ca82a306e0af3a5f20602c85c0a1de1c2f88/68747470733a2f2f7777772e7365656b6c6f676f2e6e65742f77702d636f6e74656e742f75706c6f6164732f323031342f31322f747769747465722d6c6f676f2d766563746f722d646f776e6c6f61642e6a7067",
  });

  const hls = new Hls();

  hls.attachMedia(player.video)
  hls.on(Hls.Events.MEDIA_ATTACHED, function () {
    hls.loadSource(hlsUrl)
  })

  const videobox = ref<HTMLDivElement | null>(null)
  if (getCurrentInstance()) {
    onMounted(() => {
      player.mount(unref(videobox) as HTMLDivElement)
    })
  }

const path = computed(() => useRoute().path);

const status = ref(null);

const turnOnLight = async () => {
  const toast = useToast();
  try {
    const data = await axios.get("https://jojot.singzer.cn/light/on");
    toast.success(data.data + " " + new Date().toLocaleString());
  } catch (error) {
    if (error.response) {
      toast.error(error.response.data + " " + new Date().toLocaleString());
      return;
    }
    toast.error(error.code + " " + new Date().toLocaleString());
  }
};

const turnOffLight = async () => {
  const toast = useToast();

  try {
    const data = await axios.get("https://jojot.singzer.cn/light/off");
    toast.success(data.data + " " + new Date().toLocaleString());
  } catch (error) {
    if (error.response) {
      toast.error(error.response.data + " " + new Date().toLocaleString());
      return;
    }
    toast.error(error.code + " " + new Date().toLocaleString());
  }
};

const call = async () => {
  const toast = useToast();

  try {
    const data = await axios.get("https://jojot.singzer.cn/call");
    toast.success(data.data + " " + new Date().toLocaleString());
  } catch (error) {
    if (error.response) {
      toast.error(error.response.data + " " + new Date().toLocaleString());
      return;
    }
    toast.error(error.code + " " + new Date().toLocaleString());
  }
};

const sleepMode = async () => {
  const toast = useToast();

  try {
    const data = await axios.get("https://jojot.singzer.cn/sleep");
    toast.success(data.data + " " + new Date().toLocaleString());
  } catch (error) {
    if (error.response) {
      toast.error(error.response.data + " " + new Date().toLocaleString());
      return;
    }
    toast.error(error.code + " " + new Date().toLocaleString());
  }

  await getStatus();
};

const getStatus = async () => {
  const data = await axios.get("https://jojot.singzer.cn/status");
  if (data.status === 200) {
    status.value = data.data;
  } else {
    status.value = null;
  }
};

const connWs = () => {
  const ws = new WebSocket("wss://jojot.singzer.cn/ws");
  ws.onopen = () => {
    console.log("ws open");
  };
  ws.onmessage = (e) => {
    const data = JSON.parse(e.data);
    console.log(data);
    data.type === "status" && (status.value = data.data);
  };
  ws.onclose = () => {
    console.log("ws close");
    // 重新连接
    setTimeout(() => {
      connWs();
    }, 1000);
  };
};

// const name = $ref('')
// const router = useRouter()
// const go = () => {
//   if (name)
//     router.push(`/hi/${encodeURIComponent(name)}`)
// }

const isNotSupport = ref(false);

const VideoType = ref<null | "flv" | "hls">(null);

const initVideoPlayer = () => {



  // 播放 hls
  // const video = document.querySelector("video");

  // const flvURl = "https://live.singzer.cn/live/jojo.flv";

  // VideoType.value = "hls";
  // if (Hls.isSupported()) {
  //   const hls = new Hls();
  //   hls.loadSource(hlsUrl);
  //   hls.attachMedia(video);
  //   video.play();
  //   return;
  // } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
  //   video.src = hlsUrl;
  //   return;
  // }

  // // 播放 flv
  // VideoType.value = "flv";
  // if (Flv.isSupported) {
  //   const flvPlayer = Flv.createPlayer({
  //     type: "flv",
  //     url: flvURl,
  //   });
  //   flvPlayer.attachMediaElement(video);
  //   flvPlayer.load();
  //   flvPlayer.play();
  //   return;
  // }

  // isNotSupport.value = true;
};

onMounted(async () => {
  initVideoPlayer();
  connWs();
});

onUnmounted(() => { });
</script>

<template>
  <ADialog v-model="showDialog">
    <ACard title="请JOJO吃瓜子~">
      <div py-5 px-5 flex flex-col justify-center items-center>
        <text py-1>记得备注信息哦!</text>
        <img width="256" height="256" src="/dn.jpg" />
        <ABtn class="my-3 text-sm btn" rounded-2xl @click="showDialog = false">
          关闭
        </ABtn>
      </div>
    </ACard>
  </ADialog>

  <ADialog v-model="showSleepDialog">
    <ACard title="开启睡眠模式">
      <div py-5 px-5 flex flex-col justify-center items-center>
        <text py-1>开启后将进入睡眠💤, 无法操作交互功能, 待模式结束后恢复!</text>
        <text py-1>确认开启么?</text>

        <div flex flex-row justify-center px-auto>
          <ABtn class="my-3 text-sm btn px-auto mx-10" rounded-2xl @click="showSleepDialog = false">
            取消
          </ABtn>
          <ABtn class="my-3 text-sm btn px-auto mx-10" rounded-2xl color="info" @click="sleepMode">
            确认
          </ABtn>
        </div>
      </div>
    </ACard>
  </ADialog>

  <div>
    <div text-4xl inline-block>🦜</div>
    <p>
      <a text-2xl rel="noreferrer" href="https://github.com/antfu/vitesse-lite" target="_blank">
        JOJO
      </a>
    </p>
    <p>
      <em text-xl op75>我是一只快活的傻鸟</em>
    </p>

    <p>
      <em text-sm op75>想用我的可爱治愈你~</em>
    </p>

    <div py-1 />

    <div>
      <div text-xl text-blue-5 font-bold>功能正在开发中...</div>

      <div v-if="!status" px-10 mx-auto w-sm py-1 my-1 flex flex-wrap flex-col rounded bg-green-5 text-white
        justify-center items-center>
        <div font-bold>JOJO现在出去玩啦, 等他回家吧~</div>

        <!-- <div font-bold>
          打算整一个涂鸦板的功能
        </div>

        <div text-sm>
          (利用墨水屏实现, 感谢评论区的创意~)
        </div> -->
      </div>

      <div v-if="status" px-auto mx-auto w-sm py-1 my-1 flex flex-wrap flex-col rounded bg-blue-5 text-white
        justify-center items-start>
        <div mx-auto>
          <div class="flex flex-row" justify-between>
            <div>电池电量: {{ status?.Battery.BatteryPercentage }} %</div>
          </div>
          <div class="flex flex-row">
            <div>充电状态: {{ status?.Battery.BatterISCharging ? "是" : "否" }}</div>
          </div>
          <div class="flex flex-row" justify-between>
            <div>设备温度: {{ status?.Battery.BatteryTemperature.toFixed(2) }} °C</div>
          </div>
          <div class="flex flex-row">
            <div>室内温度: {{ status?.IndoorTemperature }} °C</div>
          </div>
          <div v-if="status?.IsSleep" class="flex flex-row">
            <div>唤醒时间: {{ status?.WakeTime }}</div>
          </div>
          <div class="flex flex-row">
            <div>观看人数: {{ status?.OnlineNum }}</div>
          </div>
        </div>
      </div>

      <div v-if="status">
        <ABtn class="m-3 text-sm btn" color="info" @click="turnOnLight"> 开灯 </ABtn>

        <ABtn class="m-3 text-sm btn" color="info" @click="turnOffLight"> 关灯 </ABtn>

        <ABtn class="m-3 text-sm btn" color="success" @click="call"> 呼叫 </ABtn>

        <ABtn v-if="status && !status.IsSleep" class="m-3 text-sm btn" @click="showSleepDialog = true">
          睡眠模式
        </ABtn>
      </div>
    </div>

    <div flex flex-col justify-center items-center mx-10>
      <div v-show="status" shadow-sm>
        <div id="videobox" ref="videobox" shadow-sm class=""></div>
      </div>

      <ABtn class="my-3 text-sm btn" rounded-2xl color="warning" @click="showDialog = true">
        打赏
      </ABtn>

      <!-- <div pa-10>
        <TheInput v-model="name" placeholder="发送弹幕" autocomplete="false" @keydown.enter="go" />
      </div> -->

      <text font-bold>如果你有好的想法或者建议</text>
      <text font-bold>可以在下面评论或者联系我 (wx: oh-icepie)</text>

      <!-- <div></div> -->
    </div>

    <div mx-auto px-auto>
      <Waline :serverURL="serverURL" :path="path" dark=".dark" />
    </div>
  </div>
</template>

<style>
/* #vcomments .vcards .vcard {
    padding: 15px 20px 0 20px;
    border-radius: 10px;
    margin-bottom: 15px;
    box-shadow: 0 0 2px 1px rgba(0, 0, 0, .12);
    transition: all .3s
}

#vcomments .vcards .vcard:hover {
    box-shadow: 0 0 6px 3px rgba(0, 0, 0, .12)
}

#vcomments .vcards .vcard .vh .vcard {
    border: none;
    box-shadow: none;
} */
</style>
