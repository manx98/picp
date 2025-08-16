<script setup>
import { computed, defineEmits, defineProps, ref, shallowRef } from 'vue'
import { showWarning } from '~/utils/index.js'

const props = defineProps({
  pin: Number,
})

const emits = defineEmits(['update:pin'])
const width_value = ref(410)
const pinVal = ref(-1)
const view_width = computed(() => {
  return `${width_value.value}px`
})
const view_height_half = computed(() => {
  return `${width_value.value / 2}px`
})

const view_line_width = computed(() => {
  return `${width_value.value / 2 - 29}px`
})

const show = shallowRef(false)

function handleSelect(value) {
  const index = value.indexOf('_')
  if (index !== -1) {
    value = value.substring(index + 1)
    if (value.startsWith('gpio')) {
      pinVal.value = Number.parseInt(value.substring(4))
      return
    }
  }
  showWarning('引脚不支持')
}

function getClass(classList) {
  classList = classList.split(' ')
  const activeGpio = `gpio${pinVal.value}`
  for (let i = 0; i < classList.length; i++) {
    if (classList[i] === activeGpio) {
      classList.push('active')
      return classList
    }
  }
  return classList
}

function applyPin() {
  emits('update:pin', pinVal.value)
  show.value = false
}

function showDialog() {
  pinVal.value = props.pin
  show.value = true
}
</script>

<template>
  <el-button type="primary" @click="showDialog">
    GPIO {{ props.pin }}
  </el-button>
  <el-dialog v-model="show" title="引脚选择" :width="width_value + 30">
    <nav id="gpio">
      <div id="pinbase" />
      <ul class="bottom">
        <li :class="getClass('pin1 gnd pow3v3')">
          <a title="" @click="handleSelect('pin1_pow3v3')"><span class="default"><span class="phys">1</span> 3v3 Power</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin3 gpio i2c gpio2')">
          <a title="Wiring Pi pin 8" @click="handleSelect('pin3_gpio2')"><span class="default"><span class="phys">3</span> <span
            class="name"
          >GPIO 2</span> <small>(I2C1 SDA)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin5 gpio i2c gpio3')">
          <a title="Wiring Pi pin 9" @click="handleSelect('pin5_gpio3')"><span class="default"><span class="phys">5</span> <span
            class="name"
          >GPIO 3</span> <small>(I2C1 SCL)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin7 gpio gpio4')">
          <a title="Wiring Pi pin 7" @click="handleSelect('pin7_gpio4')"><span class="default"><span class="phys">7</span> <span
            class="name"
          >GPIO 4</span> <small>(GPCLK0)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin9 gnd')">
          <a title="" @click="handleSelect('pin9_ground')"><span class="default"><span class="phys">9</span> Ground</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin11 gpio gpio17')">
          <a title="Wiring Pi pin 0" @click="handleSelect('pin11_gpio17')"><span class="default"><span
            class="phys"
          >11</span> <span
            class="name"
          >GPIO 17</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin13 gpio gpio27')">
          <a title="Wiring Pi pin 2" @click="handleSelect('pin13_gpio27')"><span class="default"><span
            class="phys"
          >13</span> <span
            class="name"
          >GPIO 27</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin15 gpio gpio22')">
          <a title="Wiring Pi pin 3" @click="handleSelect('pin15_gpio22')"><span class="default"><span
            class="phys"
          >15</span> <span
            class="name"
          >GPIO 22</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin17 gnd pow3v3')">
          <a title="" @click="handleSelect('pin17_pow3v3')"><span class="default"><span class="phys">17</span> 3v3 Power</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin19 gpio spi gpio10')">
          <a title="Wiring Pi pin 12" @click="handleSelect('pin19_gpio10')"><span class="default"><span
            class="phys"
          >19</span> <span
            class="name"
          >GPIO 10</span> <small>(SPI0 MOSI)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin21 gpio spi gpio9')">
          <a title="Wiring Pi pin 13" @click="handleSelect('pin21_gpio9')"><span class="default"><span
            class="phys"
          >21</span> <span
            class="name"
          >GPIO 9</span> <small>(SPI0 MISO)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin23 gpio spi gpio11')">
          <a title="Wiring Pi pin 14" @click="handleSelect('pin23_gpio11')"><span class="default"><span
            class="phys"
          >23</span> <span
            class="name"
          >GPIO 11</span> <small>(SPI0 SCLK)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin25 gnd')">
          <a title="" @click="handleSelect('pin25_gnd')"><span class="default"><span class="phys">25</span> Ground</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin27 gpio i2c gpio0')">
          <a title="Wiring Pi pin 30" @click="handleSelect('pin27_gpio0')"><span class="default"><span
            class="phys"
          >27</span> <span
            class="name"
          >GPIO 0</span> <small>(EEPROM SDA)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin29 gpio gpio5')">
          <a title="Wiring Pi pin 21" @click="handleSelect('pin29_gpio5')"><span class="default"><span
            class="phys"
          >29</span> <span
            class="name"
          >GPIO 5</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin31 gpio gpio6')">
          <a title="Wiring Pi pin 22" @click="handleSelect('pin31_gpio6')"><span class="default"><span
            class="phys"
          >31</span> <span
            class="name"
          >GPIO 6</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin33 gpio gpio13')">
          <a title="Wiring Pi pin 23" @click="handleSelect('pin33_gpio13')"><span class="default"><span
            class="phys"
          >33</span> <span
            class="name"
          >GPIO 13</span> <small>(PWM1)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin35 gpio pcm gpio19')">
          <a title="Wiring Pi pin 24" @click="handleSelect('pin35_gpio19')"><span class="default"><span
            class="phys"
          >35</span> <span
            class="name"
          >GPIO 19</span> <small>(PCM FS)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin37 gpio gpio26')">
          <a title="Wiring Pi pin 25" @click="handleSelect('pin37_gpio26')"><span class="default"><span
            class="phys"
          >37</span> <span
            class="name"
          >GPIO 26</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin39 gnd')">
          <a title="" @click="handleSelect('pin39_gnd')"><span class="default"><span class="phys">39</span> Ground</span><span
            class="pin"
          /></a>
        </li>
      </ul>
      <ul class="top">
        <li :class="getClass('pin2 gnd pow5v')">
          <a title="" @click="handleSelect('pin2_pow5v')"><span class="default"><span class="phys">2</span> 5v Power</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin4 gnd pow5v')">
          <a title="" @click="handleSelect('pin4_pow5v')"><span class="default"><span class="phys">4</span> 5v Power</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin6 gnd')">
          <a title="" @click="handleSelect('pin6_gnd')"><span class="default"><span class="phys">6</span> Ground</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin8 gpio uart gpio14')">
          <a title="Wiring Pi pin 15" @click="handleSelect('pin8_gpio14')"><span class="default"><span class="phys">8</span> <span
            class="name"
          >GPIO 14</span> <small>(UART TX)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin10 gpio uart gpio15')">
          <a title="Wiring Pi pin 16" @click="handleSelect('pin10_gpio15')"><span class="default"><span
            class="phys"
          >10</span> <span
            class="name"
          >GPIO 15</span> <small>(UART RX)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin12 gpio pcm gpio18')">
          <a title="Wiring Pi pin 1" @click="handleSelect('pin12_gpio18')"><span class="default"><span
            class="phys"
          >12</span> <span
            class="name"
          >GPIO 18</span> <small>(PCM CLK)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin14 gnd')">
          <a title="" @click="handleSelect('pin14_gnd')"><span class="default"><span class="phys">14</span> Ground</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin16 gpio gpio23')">
          <a title="Wiring Pi pin 4" @click="handleSelect('pin16_gpio23')"><span class="default"><span
            class="phys"
          >16</span> <span
            class="name"
          >GPIO 23</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin18 gpio gpio24')">
          <a title="Wiring Pi pin 5" @click="handleSelect('pin18_gpio24')"><span class="default"><span
            class="phys"
          >18</span> <span
            class="name"
          >GPIO 24</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin20 gnd')">
          <a title="" @click="handleSelect('pin20_gnd')"><span class="default"><span class="phys">20</span> Ground</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin22 gpio gpio25')">
          <a title="Wiring Pi pin 6" @click="handleSelect('pin22_gpio25')"><span class="default"><span
            class="phys"
          >22</span> <span
            class="name"
          >GPIO 25</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin24 gpio spi gpio8')">
          <a title="Wiring Pi pin 10" @click="handleSelect('pin24_gpio8')"><span class="default"><span
            class="phys"
          >24</span> <span
            class="name"
          >GPIO 8</span> <small>(SPI0 CE0)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin26 gpio spi gpio7')">
          <a title="Wiring Pi pin 11" @click="handleSelect('pin26_gpio7')"><span class="default"><span
            class="phys"
          >26</span> <span
            class="name"
          >GPIO 7</span> <small>(SPI0 CE1)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin28 gpio i2c gpio1')">
          <a title="Wiring Pi pin 31" @click="handleSelect('pin28_gpio1')"><span class="default"><span
            class="phys"
          >28</span> <span
            class="name"
          >GPIO 1</span> <small>(EEPROM SCL)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin30 gnd')">
          <a title="" @click="handleSelect('pin30_gnd')"><span class="default"><span class="phys">30</span> Ground</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin32 gpio gpio12')">
          <a title="Wiring Pi pin 26" @click="handleSelect('pin32_gpio12')"><span class="default"><span
            class="phys"
          >32</span> <span
            class="name"
          >GPIO 12</span> <small>(PWM0)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin34 gnd')">
          <a title="" @click="handleSelect('pin34_gnd')"><span class="default"><span class="phys">34</span> Ground</span><span
            class="pin"
          /></a>
        </li>
        <li :class="getClass('pin36 gpio gpio16')">
          <a title="Wiring Pi pin 27" @click="handleSelect('pin36_gpio16')"><span class="default"><span
            class="phys"
          >36</span> <span
            class="name"
          >GPIO 16</span></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin38 gpio pcm gpio20')">
          <a title="Wiring Pi pin 28" @click="handleSelect('pin38_gpio20')"><span class="default"><span
            class="phys"
          >38</span> <span
            class="name"
          >GPIO 20</span> <small>(PCM DIN)</small></span><span class="pin" /></a>
        </li>
        <li :class="getClass('pin40 gpio pcm gpio21')">
          <a title="Wiring Pi pin 29" @click="handleSelect('pin40_gpio21')"><span class="default"><span
            class="phys"
          >40</span> <span
            class="name"
          >GPIO 21</span> <small>(PCM DOUT)</small></span><span class="pin" /></a>
        </li>
      </ul>
    </nav>
    <template #footer>
      <el-button type="primary" @click="applyPin">
        确认
      </el-button>
      <el-button @click="show = false">
        取消
      </el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
/*
GPIO nav
*/
nav {
  position: relative;
}

#gpio {
  width: v-bind(view_width);
  min-height: 493px;
  background: #5f8645;
}

#gpio:before {
  content: '';
  display: block;
  width: 58px;
  position: absolute;
  left: v-bind(view_line_width);
  height: 493px;
  background: #073642;
  top: 0px;
}

#gpio ul {
  position: relative;
  top: 7px;
  list-style: none;
  display: block;
  width: v-bind(view_height_half);
  float: left;
}

#gpio a {
  display: block;
  position: relative;
  font-size: 0.84em;
  line-height: 22px;
  height: 22px;
  margin-bottom: 2px;
  color: #e9e5d2;
  width: v-bind(view_height_half);
}

#gpio .phys {
  color: #073642;
  font-size: 0.8em;
  opacity: 0.8;
  position: absolute;
  left: 32px;
  text-indent: 0;
}

#gpio .pin {
  display: block;
  border: 1px solid transparent;
  border-radius: 50%;
  width: 16px;
  height: 16px;
  background: #002b36;
  position: absolute;
  right: 4px;
  top: 2px;
}

#gpio .pin:after {
  content: '';
  display: block;
  border-radius: 100%;
  background: #fdf6e3;
  position: absolute;
  left: 5px;
  top: 5px;
  width: 6px;
  height: 6px;
}

#gpio .top li {
  text-indent: 56px;
}

#gpio .top a {
  border-top-left-radius: 13px;
  border-bottom-left-radius: 13px;
}

#gpio .top .pin {
  left: 4px;
  top: 2px;
}

#gpio .bottom a {
  text-indent: 10px;
  border-top-right-radius: 13px;
  border-bottom-right-radius: 13px;
}

#gpio .bottom .overlay-ground .phys {
  padding-right: 32px;
  right: 0;
}

#gpio .bottom .phys {
  text-align: right;
  left: auto;
  right: 32px;
}

#gpio .gnd a {
  color: rgba(233, 229, 210, 0.5);
}

#gpio:hover {
  color: rgba(6, 53, 65, 0.5);
}

#gpio a:hover,
#gpio .active a {
  background: #f5f3ed;
  color: #063541;
}

#gpio li a small {
  font-size: 0.7em;
}

#gpio .overlay-pin a {
  background: #ebe6d3;
  color: #063541;
}

#gpio .overlay-pin a:hover {
  background: #f5f3ed;
  color: #063541;
}

#gpio .overlay-pin.gnd a {
  color: rgba(6, 53, 65, 0.5);
}

#gpio .overlay-power .phys {
  color: #ffffff;
  opacity: 1;
}

#gpio .overlay-power a {
  background: #073642;
  color: #ffffff;
}

#gpio .overlay-power a:hover {
  background: #268bd2;
}

#gpio .overlay-ground .phys {
  background: #073642;
  color: #ffffff;
  opacity: 1;
  position: absolute;
  top: 0px;
  width: 20px;
  height: 22px;
  border-radius: 11px;
  text-indent: 0px;
  line-height: 22px;
}

#gpio .overlay-ground a:hover .phys {
  background: #268bd2;
}

#gpio .overlay-ground span.pin {
  background: #073642;
}

#gpio ul li.hover-pin a,
#gpio .bottom li.hover-pin a {
  color: #fff;
  background: rgba(200, 0, 0, 0.6);
}

#gpio ul li.hover-pin a .phys,
#gpio .bottom li.hover-pin a .phys {
  color: #fff;
}

#gpio .pin1 a:hover,
#gpio .pin1.active a,
#gpio .pin1 .pin {
  border-radius: 0;
}

#gpio .pow3v3 .pin {
  background: #b58900;
}

#gpio .pow5v .pin {
  background: #dc322f;
}

#gpio .gpio .pin {
  background: #859900;
}

#gpio .i2c .pin {
  background: #268bd2;
}

#gpio .spi .pin {
  background: #d33682;
}

#gpio .pcm .pin {
  background: #2aa198;
}

#gpio .uart .pin {
  background: #6c71c4;
}
</style>
