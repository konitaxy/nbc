<template>
    <div class="position-relative rounded-2  text-nowrap">
        <button style="min-width: 80px;" class="btn bg-white position-relative z-3 btn-sm d-flex justify-content-between w-100 h-100" @click="handleDatePickerOpen">
            <i class="bi bi-calendar2 me-1"></i>{{dateStringShow}}
            </button>
            <el-date-picker
                v-model="dateRange"
                type="daterange"
                prefix-icon="dd"
                unlink-panels
                range-separator=""
                start-placeholder=""
                end-placeholder=""
                :shortcuts="shortcuts"
                size="small"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                ref="datePicker"
                @change="handleChangeDatePicker"
                style="width: 0;height: 0; position: absolute;left: 50%;bottom:10%;border: none;"
            />                    
    </div>
</template>
<script>

export default {
  name: 'datePicker'
}
</script>
<script setup>
import { ref, onMounted, computed, reactive,watch } from 'vue'
import {formatDateYYYYMMDD,formatDateMMDD,formatDateFormat} from '@/utils/format'

const props = defineProps({
    dateRange:{
        type:Array,
        default:[]
    }
})
const emit = defineEmits(['change'])
const datePicker = ref()

const dateRange = ref(props.dateRange)
const handleDatePickerOpen = () => {
    datePicker.value.focus()    
}
const handleChangeDatePicker =(value)=>{
    const today =  formatDateYYYYMMDD(new Date())
    var start = value[0];
    var end = value[1];
    if (start > today) {
        dateRange.value[0] = today
    }
    if (end > today) {
        dateRange.value[1] = today
    }
    
    emit('change',dateRange.value)
}
const dateStringShow = computed({
    get(){
        if(dateRange.value &&dateRange.value.length>0) {
            const today =  formatDateYYYYMMDD(new Date())
            const yesterday = formatDateYYYYMMDD(new Date().setDate(new Date().getDate()-1))
            var start = dateRange.value[0];
            var end = dateRange.value[1];
            if (start > today) {
                start = today
            }
            if (end > today) {
                end = today
            }
            console.log(start,"---",end)
            if (end == start) {
                if (end == today){
                    return "today"
                }else if (end == yesterday){
                    return "yesterday"
                }else {
                    return formatDateFormat(end,"MM/dd yyyy")
                }
            }
            const endIsToday = end == today
            start = new Date(start)
            end = new Date(end)
            const diffInMs = Math.abs(end.getTime() - start.getTime());
            const days = Math.ceil(diffInMs / (1000 * 60 * 60 * 24));
            var str = ""
            switch(days) {
                case 7:  str = endIsToday?'last week':'';break
                case 30: str = endIsToday?'last month':'';break
                case 90: str = endIsToday?'last 3 month':'';break
            }
            if (days > 5000){
                return 'All'
            }else {
                if (str == ''){
                    if(formatDateFormat(dateRange.value[0],'yyyy') == formatDateFormat(dateRange.value[1],'yyyy')){
                        return formatDateFormat(dateRange.value[0],'MM/dd') +" - " + formatDateFormat(dateRange.value[1],'MM/dd yyyy')
                    }else {
                        return formatDateFormat(dateRange.value[0],'MM/dd yyyy') +" - " + formatDateFormat(dateRange.value[1],'MM/dd yyyy')

                    }
                }else {
                    return str
                }
            }
        }else {
            return "All"
        }
       
    },
})
const shortcuts = [
{
    text: 'All',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 365*20)
      return [start,end,"all"]
    },
  },
  {
    text: 'today',
    value: () => {
      const end = new Date()
      const start = new Date()
      return [start, end,"today"]
    },
  },
  {
    text: 'yesterday',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24)
      return [start, end,'tomorrow']
    },
  },
{
    text: 'Last week',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end,'Last week']
    },
  },
  {
    text: 'Last month',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end,'Last month']
    },
  },
  {
    text: 'Last 3 month',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end,'Last 3 month']
    },
  },
]

</script>
