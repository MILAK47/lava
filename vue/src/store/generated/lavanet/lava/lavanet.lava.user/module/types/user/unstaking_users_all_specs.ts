/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import { UserStake } from "../user/user_stake";
import { SpecStakeStorage } from "../user/spec_stake_storage";

export const protobufPackage = "lavanet.lava.user";

export interface UnstakingUsersAllSpecs {
  id: number;
  unstaking: UserStake | undefined;
  specStakeStorage: SpecStakeStorage | undefined;
}

const baseUnstakingUsersAllSpecs: object = { id: 0 };

export const UnstakingUsersAllSpecs = {
  encode(
    message: UnstakingUsersAllSpecs,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    if (message.unstaking !== undefined) {
      UserStake.encode(message.unstaking, writer.uint32(18).fork()).ldelim();
    }
    if (message.specStakeStorage !== undefined) {
      SpecStakeStorage.encode(
        message.specStakeStorage,
        writer.uint32(26).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UnstakingUsersAllSpecs {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUnstakingUsersAllSpecs } as UnstakingUsersAllSpecs;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.unstaking = UserStake.decode(reader, reader.uint32());
          break;
        case 3:
          message.specStakeStorage = SpecStakeStorage.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnstakingUsersAllSpecs {
    const message = { ...baseUnstakingUsersAllSpecs } as UnstakingUsersAllSpecs;
    if (object.id !== undefined && object.id !== null) {
      message.id = Number(object.id);
    } else {
      message.id = 0;
    }
    if (object.unstaking !== undefined && object.unstaking !== null) {
      message.unstaking = UserStake.fromJSON(object.unstaking);
    } else {
      message.unstaking = undefined;
    }
    if (
      object.specStakeStorage !== undefined &&
      object.specStakeStorage !== null
    ) {
      message.specStakeStorage = SpecStakeStorage.fromJSON(
        object.specStakeStorage
      );
    } else {
      message.specStakeStorage = undefined;
    }
    return message;
  },

  toJSON(message: UnstakingUsersAllSpecs): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.unstaking !== undefined &&
      (obj.unstaking = message.unstaking
        ? UserStake.toJSON(message.unstaking)
        : undefined);
    message.specStakeStorage !== undefined &&
      (obj.specStakeStorage = message.specStakeStorage
        ? SpecStakeStorage.toJSON(message.specStakeStorage)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<UnstakingUsersAllSpecs>
  ): UnstakingUsersAllSpecs {
    const message = { ...baseUnstakingUsersAllSpecs } as UnstakingUsersAllSpecs;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = 0;
    }
    if (object.unstaking !== undefined && object.unstaking !== null) {
      message.unstaking = UserStake.fromPartial(object.unstaking);
    } else {
      message.unstaking = undefined;
    }
    if (
      object.specStakeStorage !== undefined &&
      object.specStakeStorage !== null
    ) {
      message.specStakeStorage = SpecStakeStorage.fromPartial(
        object.specStakeStorage
      );
    } else {
      message.specStakeStorage = undefined;
    }
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}