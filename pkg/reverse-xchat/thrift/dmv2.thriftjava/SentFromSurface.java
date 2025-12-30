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
@Metadata(m64929d1 = {"\u0000\u0012\n\u0002\u0018\u0002\n\u0002\u0010\u0010\n\u0000\n\u0002\u0010\b\n\u0002\b\t\b\u0086\u0081\u0002\u0018\u0000 \u000b2\b\u0012\u0004\u0012\u00020\u00000\u0001:\u0001\u000bB\u0011\b\u0002\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005R\u0010\u0010\u0002\u001a\u00020\u00038\u0006X\u0087\u0004¢\u0006\u0002\n\u0000j\u0002\b\u0006j\u0002\b\u0007j\u0002\b\bj\u0002\b\tj\u0002\b\n¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/SentFromSurface;", "", "value", "", "<init>", "(Ljava/lang/String;II)V", "CONVERSATION_SCREEN_COMPOSER", "NOTIFICATION_REPLY", "SHARE_SHEET", "PAYMENTS_SUPPORT_COMPOSER", "MESSAGE_FORWARD_SHEET", "Companion", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final class SentFromSurface {
    private static final /* synthetic */ EnumEntries $ENTRIES;
    private static final /* synthetic */ SentFromSurface[] $VALUES;

    /* renamed from: Companion, reason: from kotlin metadata */
    @InterfaceC88464a
    public static final Companion INSTANCE;

    @JvmField
    public final int value;
    public static final SentFromSurface CONVERSATION_SCREEN_COMPOSER = new SentFromSurface("CONVERSATION_SCREEN_COMPOSER", 0, 1);
    public static final SentFromSurface NOTIFICATION_REPLY = new SentFromSurface("NOTIFICATION_REPLY", 1, 2);
    public static final SentFromSurface SHARE_SHEET = new SentFromSurface("SHARE_SHEET", 2, 3);
    public static final SentFromSurface PAYMENTS_SUPPORT_COMPOSER = new SentFromSurface("PAYMENTS_SUPPORT_COMPOSER", 3, 4);
    public static final SentFromSurface MESSAGE_FORWARD_SHEET = new SentFromSurface("MESSAGE_FORWARD_SHEET", 4, 5);

    @Metadata(m64929d1 = {"\u0000\u0018\n\u0002\u0018\u0002\n\u0002\u0010\u0000\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\u0003\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0010\u0010\u0004\u001a\u0004\u0018\u00010\u00052\u0006\u0010\u0006\u001a\u00020\u0007¨\u0006\b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/SentFromSurface$Companion;", "", "<init>", "()V", "findByValue", "Lcom/x/dmv2/thriftjava/SentFromSurface;", "value", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class Companion {
        public /* synthetic */ Companion(DefaultConstructorMarker defaultConstructorMarker) {
            this();
        }

        @InterfaceC88465b
        public final SentFromSurface findByValue(int value) {
            if (value == 1) {
                return SentFromSurface.CONVERSATION_SCREEN_COMPOSER;
            }
            if (value == 2) {
                return SentFromSurface.NOTIFICATION_REPLY;
            }
            if (value == 3) {
                return SentFromSurface.SHARE_SHEET;
            }
            if (value == 4) {
                return SentFromSurface.PAYMENTS_SUPPORT_COMPOSER;
            }
            if (value != 5) {
                return null;
            }
            return SentFromSurface.MESSAGE_FORWARD_SHEET;
        }

        private Companion() {
        }
    }

    private static final /* synthetic */ SentFromSurface[] $values() {
        return new SentFromSurface[]{CONVERSATION_SCREEN_COMPOSER, NOTIFICATION_REPLY, SHARE_SHEET, PAYMENTS_SUPPORT_COMPOSER, MESSAGE_FORWARD_SHEET};
    }

    static {
        SentFromSurface[] sentFromSurfaceArr$values = $values();
        $VALUES = sentFromSurfaceArr$values;
        $ENTRIES = EnumEntriesKt.m65223a(sentFromSurfaceArr$values);
        INSTANCE = new Companion(null);
    }

    private SentFromSurface(String str, int i, int i2) {
        this.value = i2;
    }

    @InterfaceC88464a
    public static EnumEntries getEntries() {
        return $ENTRIES;
    }

    public static SentFromSurface valueOf(String str) {
        return (SentFromSurface) Enum.valueOf(SentFromSurface.class, str);
    }

    public static SentFromSurface[] values() {
        return (SentFromSurface[]) $VALUES.clone();
    }
}
