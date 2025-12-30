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
@Metadata(m64929d1 = {"\u0000\u0012\n\u0002\u0018\u0002\n\u0002\u0010\u0010\n\u0000\n\u0002\u0010\b\n\u0002\b\u0006\b\u0086\u0081\u0002\u0018\u0000 \b2\b\u0012\u0004\u0012\u00020\u00000\u0001:\u0001\bB\u0011\b\u0002\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005R\u0010\u0010\u0002\u001a\u00020\u00038\u0006X\u0087\u0004¢\u0006\u0002\n\u0000j\u0002\b\u0006j\u0002\b\u0007¨\u0006\t"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ScreenCaptureType;", "", "value", "", "<init>", "(Ljava/lang/String;II)V", "SCREENSHOT", "RECORDING", "Companion", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final class ScreenCaptureType {
    private static final /* synthetic */ EnumEntries $ENTRIES;
    private static final /* synthetic */ ScreenCaptureType[] $VALUES;

    /* renamed from: Companion, reason: from kotlin metadata */
    @InterfaceC88464a
    public static final Companion INSTANCE;

    @JvmField
    public final int value;
    public static final ScreenCaptureType SCREENSHOT = new ScreenCaptureType("SCREENSHOT", 0, 1);
    public static final ScreenCaptureType RECORDING = new ScreenCaptureType("RECORDING", 1, 2);

    @Metadata(m64929d1 = {"\u0000\u0018\n\u0002\u0018\u0002\n\u0002\u0010\u0000\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\u0003\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0010\u0010\u0004\u001a\u0004\u0018\u00010\u00052\u0006\u0010\u0006\u001a\u00020\u0007¨\u0006\b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ScreenCaptureType$Companion;", "", "<init>", "()V", "findByValue", "Lcom/x/dmv2/thriftjava/ScreenCaptureType;", "value", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class Companion {
        public /* synthetic */ Companion(DefaultConstructorMarker defaultConstructorMarker) {
            this();
        }

        @InterfaceC88465b
        public final ScreenCaptureType findByValue(int value) {
            if (value == 1) {
                return ScreenCaptureType.SCREENSHOT;
            }
            if (value != 2) {
                return null;
            }
            return ScreenCaptureType.RECORDING;
        }

        private Companion() {
        }
    }

    private static final /* synthetic */ ScreenCaptureType[] $values() {
        return new ScreenCaptureType[]{SCREENSHOT, RECORDING};
    }

    static {
        ScreenCaptureType[] screenCaptureTypeArr$values = $values();
        $VALUES = screenCaptureTypeArr$values;
        $ENTRIES = EnumEntriesKt.m65223a(screenCaptureTypeArr$values);
        INSTANCE = new Companion(null);
    }

    private ScreenCaptureType(String str, int i, int i2) {
        this.value = i2;
    }

    @InterfaceC88464a
    public static EnumEntries getEntries() {
        return $ENTRIES;
    }

    public static ScreenCaptureType valueOf(String str) {
        return (ScreenCaptureType) Enum.valueOf(ScreenCaptureType.class, str);
    }

    public static ScreenCaptureType[] values() {
        return (ScreenCaptureType[]) $VALUES.clone();
    }
}
