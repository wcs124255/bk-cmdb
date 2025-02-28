<template>
  <cmdb-sticky-layout class="sync-set-layout" ref="instancesInfo"
    v-bkloading="{ isLoading: $loading(requestIds.diffTemplateAndInstances) }">
    <template v-if="noInfo">
      <div class="no-content">
        <img src="../../assets/images/no-content.png" alt="no-content">
        <p>{{$t('无集群模板更新信息')}}</p>
        <bk-button theme="primary" @click="handleGoback">{{$t('返回')}}</bk-button>
      </div>
    </template>
    <template v-else-if="isLatestInfo">
      <div class="no-content">
        <img src="../../assets/images/latest-data.png" alt="no-content">
        <p>{{$t('最新集群模板信息')}}</p>
        <bk-button theme="primary" @click="handleGoback">{{$t('返回')}}</bk-button>
      </div>
    </template>
    <template v-else-if="diffList.length">
      <div class="main">
        <div class="title clearfix">
          <p class="fl" v-if="isSingleSync">{{$t('请确认单个实例更改信息')}}</p>
          <i18n path="请确认实例更改信息"
            tag="p"
            class="fl"
            v-else>
            <span place="count">{{setInstancesId.length}}</span>
          </i18n>
          <div class="tips fr">
            <span class="mr30">
              <i class="dot"></i>
              {{$t('新增模块')}}
            </span>
            <span class="mr30">
              <i class="dot delete"></i>
              {{$t('删除模块')}}
            </span>
            <bk-checkbox class="expand-all"
              v-if="!isSingleSync"
              v-model="expandAll"
              @change="handleExpandAll">
              {{$t('全部展开')}}
            </bk-checkbox>
          </div>
        </div>
        <div class="instance-list">
          <set-instance class="instance-item"
            ref="setInstance"
            v-for="(instance, index) in diffList"
            :key="instance.bk_set_id"
            :instance="instance"
            :icon-close="diffList.length > 1"
            :expand="index === 0"
            :module-host-count="moduleHostCount"
            @close="handleCloseInstance">
          </set-instance>
        </div>
      </div>
    </template>
    <div class="options" slot="footer" slot-scope="{ sticky }"
      v-if="!noInfo && !isLatestInfo && diffList.length"
      :class="{ 'is-sticky': sticky }">
      <cmdb-auth :auth="{ type: $OPERATION.U_TOPO, relation: [bizId] }">
        <template slot-scope="{ disabled }">
          <bk-button v-if="canSyncStatus()"
            class="mr10"
            theme="primary"
            :loading="$loading(requestIds.syncTemplateToInstances)"
            :disabled="disabled"
            @click="handleConfirmSync">
            {{$t('确认同步')}}
          </bk-button>
          <span class="text-btn mr10" v-else v-bk-tooltips="$t('请先删除不可同步的实例')">{{$t('确认同步')}}</span>
        </template>
      </cmdb-auth>
      <bk-button class="mr10" @click="handleGoback">{{$t('取消')}}</bk-button>
    </div>
  </cmdb-sticky-layout>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
  import setInstance from './set-instance'
  export default {
    components: {
      setInstance
    },
    data() {
      return {
        expandAll: false,
        diffList: [],
        noInfo: false,
        isLatestInfo: false,
        templateName: '',
        moduleHostCount: {},
        requestIds: {
          diffTemplateAndInstances: Symbol(),
          syncTemplateToInstances: Symbol()
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      setTemplateId() {
        return this.$route.params.setTemplateId
      },
      setInstancesId() {
        const id = `${this.bizId}_${this.setTemplateId}`
        let { syncIdMap } = this.$store.state.setFeatures
        const sessionSyncIdMap = sessionStorage.getItem('setSyncIdMap')
        if (!Object.keys(syncIdMap).length && sessionSyncIdMap) {
          syncIdMap = JSON.parse(sessionSyncIdMap)
          this.$store.commit('setFeatures/resetSyncIdMap', syncIdMap)
        }
        return syncIdMap[id] || []
      },
      isSingleSync() {
        return !(this.setInstancesId.length > 1)
      }
    },
    async created() {
      this.getSetTemplateInfo()
      this.getDiffData()
    },
    methods: {
      canSyncStatus() {
        let status = true
        this.$refs.setInstance.forEach((instance) => {
          if (!instance.canSyncStatus && status) {
            status = false
          }
        })
        return status
      },
      async getSetTemplateInfo() {
        try {
          const info = await this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
            bizId: this.bizId,
            setTemplateId: this.setTemplateId
          })
          this.templateName = info.name
        } catch (e) {
          console.error(e)
        }
      },
      async getDiffData() {
        try {
          if (!this.setInstancesId.length) {
            this.diffList = []
            this.noInfo = true
            return
          }
          const data = await this.$store.dispatch('setSync/diffTemplateAndInstances', {
            bizId: this.bizId,
            setTemplateId: this.setTemplateId,
            params: {
              bk_set_ids: this.setInstancesId
            },
            config: {
              requestId: this.requestIds.diffTemplateAndInstances
            }
          })
          this.moduleHostCount = data.module_host_count || {}
          const list = data.difference || []
          this.diffList = list.sort((prev, next) => {
            const res = module => (module.diff_type === 'remove' && this.moduleHostCount[module.bk_module_id] > 0)
            const prevEixstHostList = prev.module_diffs.filter(module => res(module))
            const nextEixstHostList = next.module_diffs.filter(module => res(module))
            if ((prevEixstHostList.length && nextEixstHostList.length)
              || (!prevEixstHostList.length && !nextEixstHostList.length)) {
              return 0
            } if (prevEixstHostList.length) {
              return -1
            }
            return 1
          })
          const changeList = this.diffList.filter((set) => {
            const moduleDiffs = set.module_diffs
            return moduleDiffs && moduleDiffs.filter(module => module.diff_type !== 'unchanged').length
          })
          this.isLatestInfo = !changeList.length
          this.noInfo = false
        } catch (e) {
          console.error(e)
          this.noInfo = true
          this.moduleHostCount = {}
        }
      },
      async handleConfirmSync() {
        try {
          await this.$store.dispatch('setSync/syncTemplateToInstances', {
            bizId: this.bizId,
            setTemplateId: this.setTemplateId,
            params: {
              bk_set_ids: this.setInstancesId
            },
            config: {
              requestId: this.requestIds.syncTemplateToInstances
            }
          })
          this.$success(this.$t('提交同步成功'))
          this.$routerActions.redirect({
            name: 'setTemplateConfig',
            params: {
              templateId: this.setTemplateId,
              mode: 'view'
            },
            query: {
              tab: 'instance'
            }
          })
        } catch (e) {
          console.error(e)
        }
      },
      handleExpandAll(expand) {
        this.$refs.setInstance.forEach((instance) => {
          instance.localExpand = expand
        })
      },
      handleCloseInstance(id) {
        this.$store.commit('setFeatures/deleteInstancesId', {
          id: `${this.bizId}_${this.setTemplateId}`,
          deleteId: id
        })
        this.diffList = this.diffList.filter(instance => instance.bk_set_id !== id)
      },
      handleGoback() {
        const { moduleId } = this.$route.params
        if (moduleId) {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_HOST_AND_SERVICE,
            query: {
              node: `set-${moduleId}`
            }
          })
        } else {
          this.$routerActions.redirect({
            name: 'setTemplateConfig',
            params: {
              templateId: this.setTemplateId,
              mode: 'view'
            },
            query: {
              tab: 'instance'
            }
          })
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .sync-set-layout {
        padding: 15px 0 0 0;
        @include scrollbar-y;
    }
    .no-content {
        position: absolute;
        top: 50%;
        left: 50%;
        font-size: 16px;
        color: #63656e;
        text-align: center;
        transform: translate(-50%, -50%);
        img {
            width: 130px;
        }
        p {
            padding: 20px 0 30px;
        }
    }
    .main {
        padding: 0 20px;
    }
    .title {
        font-size: 14px;
        color: #63656E;
    }
    .tips {
        display: flex;
        align-items: center;
        font-size: 12px;
        .dot {
            display: inline-block;
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #2DCB56;
            margin-right: 2px;
            &.delete {
                background-color: #C4C6CC;
            }
        }
    }
    .expand-all {
        color: #888991;
        /deep/ .bk-checkbox-text {
            font-size: 12px;
        }
    }
    .instance-list {
        padding: 20px 0 0;
        .instance-item {
            margin-bottom: 10px;
        }
    }
    .options {
        display: flex;
        align-items: center;
        padding: 10px 20px;
        border-top: 1px solid transparent;
        background-color: #fafbfd;
        &.is-sticky {
            background-color: #ffffff;
            border-top-color: #dcdee5;
            z-index: 100;
        }
        > span {
            color: #979BA5;
            font-size: 14px;
        }
        .text-btn {
            @include inlineBlock;
            height: 32px;
            line-height: 32px;
            outline: none;
            padding: 0 15px;
            text-align: center;
            font-size: 14px;
            background-color: #DCDEE5;
            border-radius: 2px;
            color: #FFFFFF;
            min-width: 68px;
            cursor: not-allowed;
        }
    }
</style>

<style lang="scss">
    .set-confirm-sync {
        .bk-dialog-content {
            width: 440px !important;
        }
        .bk-dialog-type-body {
            padding: 2px 24px 5px !important;
        }
        .bk-dialog-type-sub-header {
            padding: 3px 40px 24px !important;
            .header {
                white-space: unset !important;
                text-align: left !important;
            }
        }
        .bk-dialog-footer {
            padding-bottom: 32px !important;
        }
    }
</style>
