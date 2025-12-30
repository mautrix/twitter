package com.x.dmv2.thriftjava;

import kotlin.Metadata;
import kotlin.enums.EnumEntries;
import kotlin.enums.EnumEntriesKt;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

/* JADX WARN: Failed to restore enum class, 'enum' modifier and super class removed */
/* JADX WARN: Unknown enum class pattern. Please report as an issue! */
@Metadata(m64929d1 = {"\u0000\u0012\n\u0002\u0018\u0002\n\u0002\u0010\u0010\n\u0000\n\u0002\u0010\b\n\u0002\b\f\b\u0086\u0081\u0002\u0018\u0000 \u000e2\b\u0012\u0004\u0012\u00020\u00000\u0001:\u0001\u000eB\u0011\b\u0002\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005R\u0010\u0010\u0002\u001a\u00020\u00038\u0006X\u0087\u0004¢\u0006\u0002\n\u0000j\u0002\b\u0006j\u0002\b\u0007j\u0002\b\bj\u0002\b\tj\u0002\b\nj\u0002\b\u000bj\u0002\b\fj\u0002\b\r¨\u0006\u000f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/FailureType;", "", "value", "", "<init>", "(Ljava/lang/String;II)V", "EMPTY_DETAIL", "INTERNAL_ERROR", "CONTENTS_TOO_LARGE", "TOO_MANY_MESSAGES", "INVALID_SENDER_SIGNATURE", "NON_LATEST_CKEY_VERSION", "RECIPIENT_HAS_NOT_TRUSTED_CONVERSATION", "RECIPIENT_KEY_HAS_CHANGED", "Companion", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final class FailureType {
    private static final /* synthetic */ EnumEntries $ENTRIES;
    private static final /* synthetic */ FailureType[] $VALUES;

    /* renamed from: Companion, reason: from kotlin metadata */
    @InterfaceC88464a
    public static final Companion INSTANCE;

    @JvmField
    public final int value;
    public static final FailureType EMPTY_DETAIL = new FailureType("EMPTY_DETAIL", 0, 1);
    public static final FailureType INTERNAL_ERROR = new FailureType("INTERNAL_ERROR", 1, 2);
    public static final FailureType CONTENTS_TOO_LARGE = new FailureType("CONTENTS_TOO_LARGE", 2, 3);
    public static final FailureType TOO_MANY_MESSAGES = new FailureType("TOO_MANY_MESSAGES", 3, 4);
    public static final FailureType INVALID_SENDER_SIGNATURE = new FailureType("INVALID_SENDER_SIGNATURE", 4, 5);
    public static final FailureType NON_LATEST_CKEY_VERSION = new FailureType("NON_LATEST_CKEY_VERSION", 5, 6);
    public static final FailureType RECIPIENT_HAS_NOT_TRUSTED_CONVERSATION = new FailureType("RECIPIENT_HAS_NOT_TRUSTED_CONVERSATION", 6, 7);
    public static final FailureType RECIPIENT_KEY_HAS_CHANGED = new FailureType("RECIPIENT_KEY_HAS_CHANGED", 7, 8);

    @Metadata(m64929d1 = {"\u0000\u0018\n\u0002\u0018\u0002\n\u0002\u0010\u0000\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\u0003\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0010\u0010\u0004\u001a\u0004\u0018\u00010\u00052\u0006\u0010\u0006\u001a\u00020\u0007¨\u0006\b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/FailureType$Companion;", "", "<init>", "()V", "findByValue", "Lcom/x/dmv2/thriftjava/FailureType;", "value", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class Companion {
        public /* synthetic */ Companion(DefaultConstructorMarker defaultConstructorMarker) {
            this();
        }

        @InterfaceC88465b
        public final FailureType findByValue(int value) {
            switch (value) {
                case 1:
                    return FailureType.EMPTY_DETAIL;
                case 2:
                    return FailureType.INTERNAL_ERROR;
                case 3:
                    return FailureType.CONTENTS_TOO_LARGE;
                case 4:
                    return FailureType.TOO_MANY_MESSAGES;
                case 5:
                    return FailureType.INVALID_SENDER_SIGNATURE;
                case 6:
                    return FailureType.NON_LATEST_CKEY_VERSION;
                case 7:
                    return FailureType.RECIPIENT_HAS_NOT_TRUSTED_CONVERSATION;
                case 8:
                    return FailureType.RECIPIENT_KEY_HAS_CHANGED;
                default:
                    return null;
            }
        }

        private Companion() {
        }
    }

    private static final /* synthetic */ FailureType[] $values() {
        return new FailureType[]{EMPTY_DETAIL, INTERNAL_ERROR, CONTENTS_TOO_LARGE, TOO_MANY_MESSAGES, INVALID_SENDER_SIGNATURE, NON_LATEST_CKEY_VERSION, RECIPIENT_HAS_NOT_TRUSTED_CONVERSATION, RECIPIENT_KEY_HAS_CHANGED};
    }

    static {
        FailureType[] failureTypeArr$values = $values();
        $VALUES = failureTypeArr$values;
        $ENTRIES = EnumEntriesKt.m65223a(failureTypeArr$values);
        INSTANCE = new Companion(null);
    }

    private FailureType(String str, int i, int i2) {
        this.value = i2;
    }

    @InterfaceC88464a
    public static EnumEntries getEntries() {
        return $ENTRIES;
    }

    public static FailureType valueOf(String str) {
        return (FailureType) Enum.valueOf(FailureType.class, str);
    }

    public static FailureType[] values() {
        return (FailureType[]) $VALUES.clone();
    }
}
