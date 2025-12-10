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
@Metadata(m64929d1 = {"\u0000\u0012\n\u0002\u0018\u0002\n\u0002\u0010\u0010\n\u0000\n\u0002\u0010\b\n\u0002\b\n\b\u0086\u0081\u0002\u0018\u0000 \f2\b\u0012\u0004\u0012\u00020\u00000\u0001:\u0001\fB\u0011\b\u0002\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005R\u0010\u0010\u0002\u001a\u00020\u00038\u0006X\u0087\u0004¢\u0006\u0002\n\u0000j\u0002\b\u0006j\u0002\b\u0007j\u0002\b\bj\u0002\b\tj\u0002\b\nj\u0002\b\u000b¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MediaType;", "", "value", "", "<init>", "(Ljava/lang/String;II)V", "IMAGE", "GIF", "VIDEO", "AUDIO", "FILE", "SVG", "Companion", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final class MediaType {
    private static final /* synthetic */ EnumEntries $ENTRIES;
    private static final /* synthetic */ MediaType[] $VALUES;

    /* renamed from: Companion, reason: from kotlin metadata */
    @InterfaceC88464a
    public static final Companion INSTANCE;

    @JvmField
    public final int value;
    public static final MediaType IMAGE = new MediaType("IMAGE", 0, 1);
    public static final MediaType GIF = new MediaType("GIF", 1, 2);
    public static final MediaType VIDEO = new MediaType("VIDEO", 2, 3);
    public static final MediaType AUDIO = new MediaType("AUDIO", 3, 4);
    public static final MediaType FILE = new MediaType("FILE", 4, 5);
    public static final MediaType SVG = new MediaType("SVG", 5, 6);

    @Metadata(m64929d1 = {"\u0000\u0018\n\u0002\u0018\u0002\n\u0002\u0010\u0000\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\u0003\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0010\u0010\u0004\u001a\u0004\u0018\u00010\u00052\u0006\u0010\u0006\u001a\u00020\u0007¨\u0006\b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MediaType$Companion;", "", "<init>", "()V", "findByValue", "Lcom/x/dmv2/thriftjava/MediaType;", "value", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class Companion {
        public /* synthetic */ Companion(DefaultConstructorMarker defaultConstructorMarker) {
            this();
        }

        @InterfaceC88465b
        public final MediaType findByValue(int value) {
            switch (value) {
                case 1:
                    return MediaType.IMAGE;
                case 2:
                    return MediaType.GIF;
                case 3:
                    return MediaType.VIDEO;
                case 4:
                    return MediaType.AUDIO;
                case 5:
                    return MediaType.FILE;
                case 6:
                    return MediaType.SVG;
                default:
                    return null;
            }
        }

        private Companion() {
        }
    }

    private static final /* synthetic */ MediaType[] $values() {
        return new MediaType[]{IMAGE, GIF, VIDEO, AUDIO, FILE, SVG};
    }

    static {
        MediaType[] mediaTypeArr$values = $values();
        $VALUES = mediaTypeArr$values;
        $ENTRIES = EnumEntriesKt.m65223a(mediaTypeArr$values);
        INSTANCE = new Companion(null);
    }

    private MediaType(String str, int i, int i2) {
        this.value = i2;
    }

    @InterfaceC88464a
    public static EnumEntries getEntries() {
        return $ENTRIES;
    }

    public static MediaType valueOf(String str) {
        return (MediaType) Enum.valueOf(MediaType.class, str);
    }

    public static MediaType[] values() {
        return (MediaType[]) $VALUES.clone();
    }
}