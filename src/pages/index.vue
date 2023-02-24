<script setup lang="ts">
import NPlayer from "nplayer";
import Danmaku from "@nplayer/danmaku";
import { BulletOption } from "@nplayer/danmaku/dist/src/ts/danmaku/bullet";
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

const serverURL = "https://waline.singzer.cn";

const hlsUrl = "http://fit-office.singzer.cn:7002/live/jojo.m3u8";

const miaoHlsUrl = "http://fit-office.singzer.cn:7002/live/miao.m3u8";

const wsUrl = "wss://jojot.singzer.cn/ws";

let ws = new WebSocket(wsUrl);

const initWs = () => {
  ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log("ws open");
  };

  ws.onmessage = (e) => {
    const data = JSON.parse(e.data);
    console.log(data);
    switch (data.type) {
      case "status":
        status.value = data.data;
        break;
      case "danmaku":
        if (!data.data.isMe) {
          return;
        }
        player.danmaku.send({ ...data.data, time: player.currentTime, isMe: false });
        break;
      default:
        break;
    }
  };

  ws.onclose = () => {
    console.log("ws close");
    // 重新连接
    setTimeout(() => {
      initWs();
    }, 1000);
  };
};

const danmakuOptions = {
  items: [{ time: 1, text: "弹幕功能可以使用啦~" }],
  autoInsert: true,
};

const jojoPlayer = new NPlayer({
  // settings: [],
  controls: [
    ["play", "spacer", "web-fullscreen", "fullscreen"],
    ["progress"],
    ["volume"],
  ],
  bpControls: {
    650: [
      ["play", "progress", "web-fullscreen", "fullscreen"],
      ["danmaku-send", "danmaku-settings"],
      ["volume"],
    ],
  },
  live: true,
  plugins: [new Danmaku(danmakuOptions)],
  // poster:
  //   "https://photo7n.gracg.com/uploadfile/photo/2017/9/pic_se9hmr4k5qsjl81soeav5nrfw74i60z6.jpg?imageMogr2/auto-orient/thumbnail/1200x/blur/1x0/quality/98",
});


jojoPlayer.on("DanmakuSend", (opts: BulletOption) => {
  if (opts.isMe) {
    // 发送弹幕
    ws.send(
      JSON.stringify({
        type: "danmaku",
        data: {
          ...opts,
        },
      })
    );
  }
  console.log(opts);
});

if (Hls.isSupported()) {
  const hls = new Hls();
  hls.loadSource(hlsUrl);
  hls.attachMedia(jojoPlayer.video);
} else if (jojoPlayer.video.canPlayType("application/vnd.apple.mpegurl")) {
  jojoPlayer.video.src = hlsUrl;
}

const videobox2 = ref<HTMLDivElement | null>(null);
if (getCurrentInstance()) {
  onMounted(() => {
    jojoPlayer.mount(unref(videobox2) as HTMLDivElement);
  });
}

const miaoPlayer = new NPlayer({
  // settings: [],
  controls: [
    ["play", "spacer", "web-fullscreen", "fullscreen"],
    ["progress"],
    ["volume"],
  ],
  bpControls: {
    650: [
      ["play", "progress", "web-fullscreen", "fullscreen"],
      ["volume"],
    ],
  },
  live: true,
  // poster:
  //   "https://photo7n.gracg.com/uploadfile/photo/2017/9/pic_se9hmr4k5qsjl81soeav5nrfw74i60z6.jpg?imageMogr2/auto-orient/thumbnail/1200x/blur/1x0/quality/98",
});


if (Hls.isSupported()) {
  const hls = new Hls();
  hls.loadSource(miaoHlsUrl);
  hls.attachMedia(miaoPlayer.video);
} else if (miaoPlayer.video.canPlayType("application/vnd.apple.mpegurl")) {
  miaoPlayer.video.src = miaoHlsUrl;
}

const videobox1 = ref<HTMLDivElement | null>(null);
if (getCurrentInstance()) {
  onMounted(() => {
    miaoPlayer.mount(unref(videobox1) as HTMLDivElement);
  });
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

  showSleepDialog.value = false;

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

onMounted(async () => {
  initWs();
});

onUnmounted(() => {});
</script>

<template w-screen>
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
          <ABtn
            class="my-3 text-sm btn px-auto mx-10"
            rounded-2xl
            @click="showSleepDialog = false"
          >
            取消
          </ABtn>
          <ABtn
            class="my-3 text-sm btn px-auto mx-10"
            rounded-2xl
            color="info"
            @click="sleepMode"
          >
            确认
          </ABtn>
        </div>
      </div>
    </ACard>
  </ADialog>

  <div>
    <div text-4xl inline-block>🦜</div>
    <p>
      <a
        text-2xl
        rel="noreferrer"
        href="https://github.com/icepie"
        target="_blank"
      >
        JOJO
      </a>
    </p>

    <p>
      <em text-xl op75>和神出鬼没的小猫</em>
    </p>

    <!-- <p>
      <em text-xl op75>我是一只快活的傻鸟</em>
    </p> -->
<!--
    <p>
      <em text-sm op75>想用我的可爱治愈你~</em>
    </p> -->

    <div py-1 />

    <div>
      <div text-xl text-blue-5 font-bold>功能正在开发中...</div>

      <div
        v-if="!status"
        w-auto
        md:w-md
        mx-auto
        px-auto
        py-1
        my-1
        flex
        flex-wrap
        flex-col
        rounded
        bg-green-5
        text-white
        justify-center
        items-center
      >
        <div font-bold>JOJO现在出去玩啦, 等他回家吧~</div>

        <!-- <div font-bold>
          打算整一个涂鸦板的功能
        </div>

        <div text-sm>
          (利用墨水屏实现, 感谢评论区的创意~)
        </div> -->
      </div>

      <div
        v-if="status"
        w-auto
        md:w-md
        mx-auto
        px-auto
        py-1
        my-1
        flex
        flex-wrap
        flex-col
        rounded
        bg-blue-5
        text-white
        justify-center
        items-start
      >
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
          <div v-if="status?.IndoorTemperature > 0" class="flex flex-row">
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
        <!-- <ABtn class="m-3 text-sm btn" color="info" @click="turnOnLight"> 开灯 </ABtn>

        <ABtn class="m-3 text-sm btn" color="info" @click="turnOffLight"> 关灯 </ABtn> -->

        <ABtn class="m-3 text-sm btn" color="success" @click="call"> 呼叫 </ABtn>

        <ABtn
          v-if="status && !status.IsSleep"
          class="m-3 text-sm btn"
          @click="showSleepDialog = true"
        >
          睡眠模式
        </ABtn>
      </div>
    </div>

    <div flex flex-col justify-center items-center px-auto mx-auto>
      <div v-show="status">
        <div class="m-3" id="videobox1" ref="videobox1" shadow-sm w-auto md:w-md></div>
        <div  class="m-3" id="videobox2" ref="videobox2" shadow-sm w-auto md:w-md></div>
      </div>

      <iframe
        frameborder="yes"
        border="0"
        marginwidth="0"
        marginheight="0"
        width="330"
        height="110"
        src="//music.163.com/outchain/player?type=0&id=7821822508&auto=0&height=90"
      ></iframe>

      <ABtn
        class="my-3 text-sm btn"
        rounded-2xl
        color="warning"
        @click="showDialog = true"
      >
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
