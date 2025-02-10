<template>
  <div class="flex size-full justify-center p-8">
    <!-- Left Area -->
    <div class="mr-4 flex h-full w-72 flex-col space-y-4">
      <div class="flex justify-between">
        <div
          v-if="musicInfo.Picture == ''"
          class="--wails-drop-target t-card--bordered flex size-44 items-center justify-center rounded-md bg-(--td-bg-color-container)"
          :style="{ '--wails-drop-target': 'drop' }"
        >
          <div v-if="musicInfo.name" class="text-center text-base">
            {{ musicInfo.name }}
          </div>
          <t-space v-else direction="vertical" align="center">
            <div class="text-base">Drag file to this area</div>
            <UploadIcon class="!size-12" />
          </t-space>
        </div>
        <t-image v-else :src="musicInfo.Picture" class="size-44" />

        <t-space direction="vertical">
          <t-button :block="true" @click="openMusicFile">Select File</t-button>
          <t-button :block="true" @click="reset">Reset</t-button>
          <t-divider class="!my-1" />
          <t-button :block="true" @click="saveMusicFile">Save</t-button>
        </t-space>
      </div>

      <div
        class="t-card--bordered flex min-h-0 flex-1 flex-col rounded-md bg-(--td-bg-color-container)"
      >
        <div
          class="t-card__title--bordered px-2 py-3 text-base font-semibold text-(--td-brand-color-9)"
        >
          File Info
        </div>
        <div class="min-h-0 flex-1 overflow-auto">
          <div class="flex flex-col space-y-4 p-2">
            <p>
              <span class="text-base font-medium">Name：</span
              >{{ musicInfo.name }}
            </p>
            <p>
              <span class="text-base font-medium">Size：</span
              >{{ musicInfo.size }}
            </p>
            <p>
              <span class="text-base font-medium">SampleRate：</span
              >{{ musicInfo.sampleRate }}
            </p>
            <p>
              <span class="text-base font-medium">Path：</span>{{ musicPath }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Area -->
    <div
      class="t-card--bordered min-w-0 flex-1 overflow-auto rounded-md bg-(--td-bg-color-container)"
    >
      <div class="p-4">
        <t-form :label-width="150" class="!text-sm">
          <template v-for="commentKey in commentKeys" :key="commentKey">
            <t-form-item v-if="commentKey !== 'LYRICS'" :label="commentKey">
              <t-input
                v-model="musicInfo.comments[commentKey]"
                placeholder=""
              />
            </t-form-item>

            <t-form-item v-else :label="commentKey">
              <t-textarea
                v-model="musicInfo.comments[commentKey]"
                :autosize="{ minRows: 5 }"
              />
            </t-form-item>
          </template>
        </t-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, reactive } from "vue";
import { UploadIcon } from "tdesign-icons-vue-next";
import { OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime";
import {
  OpenMusicFile,
  ParseMusicFile,
  SaveMusicFile,
  CheckMusicType,
} from "../../wailsjs/go/main/App";
import { MessagePlugin } from "tdesign-vue-next";

const musicPath = ref("");
const musicInfo = reactive({
  name: computed(() => musicPath.value.split("/").pop()),
  size: "",
  sampleRate: "",
  comments: {} as Record<string, string>,
  Picture: "",
});

const commentKeys: string[] = [
  "TITLE",
  "VERSION",
  "ALBUM",
  "ALBUMARTIST",
  "ARTIST",
  "LYRICIST",
  "COMPOSER",
  "PERFORMER",
  "TRACKNUMBER",
  "COPYRIGHT",
  "LICENSE",
  "ORGANIZATION",
  "DESCRIPTION",
  "GENRE",
  "DATE",
  "LOCATION",
  "CONTACT",
  "ISRC",
  "LYRICS",
];

onMounted(() => {
  OnFileDrop((x: number, y: number, paths: string[]) => {
    if (paths.length > 0 && musicPath.value !== paths[0]) {
      musicPath.value = paths[0];
      parseMusicFile();
    }
  }, true);
});

onUnmounted(() => {
  OnFileDropOff();
});

const openMusicFile = async () => {
  const openResult = await OpenMusicFile();
  if (openResult) {
    musicPath.value = openResult;
    await parseMusicFile();
  }
};

const reset = () => {
  musicPath.value = "";
  musicInfo.size = "";
  musicInfo.sampleRate = "";
  musicInfo.comments = {};
  musicInfo.Picture = "";
};

const parseMusicFile = async () => {
  const checkResult = await CheckMusicType(musicPath.value);
  if (!checkResult) {
    await MessagePlugin.error("Only FLAC file is allowed.");
    reset();
    return;
  }
  const parseResult = await ParseMusicFile(musicPath.value);
  musicInfo.size = parseResult.Size;
  musicInfo.sampleRate = parseResult.SampleRate;
  musicInfo.comments = parseResult.Comments;
  musicInfo.Picture = parseResult.Picture;
};

const saveMusicFile = async () => {
  const saveResult = await SaveMusicFile(musicPath.value, musicInfo.comments);
  if (saveResult) await MessagePlugin.success("Save successful.");
  else await MessagePlugin.error("Save failed.");
};
</script>
