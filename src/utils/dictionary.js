import { useDictionaryStore } from '@/pinia/modules/dictionary'
export const getDict = async(type) => {
  const dictionaryStore = useDictionaryStore()
  await dictionaryStore.getDictionary(type)
  return dictionaryStore.dictionaryMap[type]
}
