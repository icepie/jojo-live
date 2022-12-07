<script setup lang="ts">

import Hls from 'hls.js'
import Flv from 'flv.js'

import Valine from "valine";

import axios from 'axios'

defineOptions({
  name: 'IndexPage',
})

const turnOnLight = () => {
  axios.get('https://jojot.singzer.cn/light/on')
}

const turnOffLight = () => {
  axios.get('https://jojot.singzer.cn/light/off')
}

const call = () => {
  axios.get('https://jojot.singzer.cn/call')
}

const getStatus = () => {
  axios.get('https://jojot.singzer.cn/status')
}

// const name = $ref('')
// const router = useRouter()
// const go = () => {
//   if (name)
//     router.push(`/hi/${encodeURIComponent(name)}`)
// }

const isNotSupport = ref(false)

const VideoType = ref<null | 'flv' | 'hls'>(null)

const initVideoPlayer = (() => {
  // 播放 hls
  const video = document.querySelector('video')
  const hlsUrl = 'https://jojos.singzer.cn/live/movie.m3u8'
  const flvURl = 'https://jojo.singzer.cn/live/movie.flv'

  VideoType.value = 'hls'
  if (Hls.isSupported()) {
    const hls = new Hls()
    hls.loadSource(hlsUrl)
    hls.attachMedia(video)
    video.play()
    return
  } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
    video.src = hlsUrl;
    return
  }
  // 播放 flv
  VideoType.value = 'flv'
  if (Flv.isSupported) {
    const flvPlayer = Flv.createPlayer({
      type: 'flv',
      url: flvURl
    })
    flvPlayer.attachMediaElement(video)
    flvPlayer.load()
    flvPlayer.play()
    return
  }

  isNotSupport.value = true
})

const initValineComment = (() => {
  new Valine({
    el: "#vcomments",
    appId: "rUxninURp0tKz3PUoEKVB4Jw-gzGzoHsz",
    appKey: "vuh6OflApSNqG84hj0kHmYbY",
    placeholder: '欢迎留言',
    visitor: true,
    avatar: 'monsterid',
    recordIP: true,
    requiredFields: ['nick', 'mail'],
    lang: 'zh-cn',
  })
})

onMounted(() => {
  initVideoPlayer()
  initValineComment()
})


</script>

<template>
  <div>
    <div text-4xl inline-block>🦜</div>
    <p>
      <a text-2xl rel="noreferrer" href="https://github.com/antfu/vitesse-lite" target="_blank">
        JOJO
      </a>
    </p>
    <p>
      <em text-xl op75>我是一只快活的傻鸟~</em>
    </p>

    <div py-1 />

    <div>
      <div text-xl text-blue-5 font-bold>功能正在开发中...</div>
      <div>
        <button class="m-3 text-sm btn" @click="turnOnLight">
          开灯
        </button>

        <button class="m-3 text-sm btn" @click="turnOffLight" >
          关灯
        </button>

        <button class="m-3 text-sm btn" @click="call">
          呼叫
        </button>
      </div>

    </div>

    <div flex flex-col justify-center items-center>

      <div shadow-sm>
        <video rounded shadow controls autoplay id="video" width="360" height="640"></video>
      </div>


      <!-- <div pa-10>
        <TheInput v-model="name" placeholder="发送弹幕" autocomplete="false" @keydown.enter="go" />
      </div> -->

      <div  py-4/>
      <text font-bold >如果你有好的想法或者建议</text>
      <text font-bold >可以在下面评论或者联系我 (wx: oh-icepie)</text>

      <div>

      </div>

      <div my-5 >
        <div id="vcomments"></div>
      </div>

      <span text-gray text-sm font-bold id="/" class="leancloud_visitors" data-flag-title="Your Article Title">
        <text class="post-meta-item-text">访问量: </text>
        <text  class="leancloud-visitors-count">1000000</text> 次
     </span>

    </div>


    <!-- <div bg-blue>
      <button class="m-3 text-sm btn" :disabled="!name" @click="go">
        Go
      </button>
    </div> -->

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
