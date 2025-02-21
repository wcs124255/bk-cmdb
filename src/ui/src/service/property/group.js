import http from '@/api'
import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'

export const find = async (options, config = {}) => {
  try {
    const params = {}
    options.bk_biz_id && (params.bk_biz_id = options.bk_biz_id)
    const groups = await http.post(`find/objectattgroup/object/${options.bk_obj_id}`, params, config)
    const bizGroups = groups.filter(group => !!group.bk_biz_id)
      .sort((previous, next) => previous.bk_group_index - next.bk_group_index)
    const globalGroups = groups.filter(group => !group.bk_biz_id)
      .sort((previous, next) => previous.bk_group_index - next.bk_group_index)
    return [...globalGroups, ...bizGroups]
  } catch (error) {
    console.error(error)
    return []
  }
}

export const findBizSet = config => find({ bk_obj_id: BUILTIN_MODELS.BUSINESS_SET }, config)

export default {
  find,
  findBizSet
}
