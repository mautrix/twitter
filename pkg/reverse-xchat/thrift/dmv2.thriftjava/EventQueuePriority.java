package com.x.dmv2.thriftjava;

import android.gov.nist.javax.sip.parser.TokenNames;
import kotlin.Metadata;
import kotlin.enums.EnumEntries;
import kotlin.enums.EnumEntriesKt;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

/* JADX WARN: Failed to restore enum class, 'enum' modifier and super class removed */
/* JADX WARN: Unknown enum class pattern. Please report as an issue! */
@Metadata(m64929d1 = {"\u0000\u0012\n\u0002\u0018\u0002\n\u0002\u0010\u0010\n\u0000\n\u0002\u0010\b\n\u0002\b\t\b\u0086\u0081\u0002\u0018\u0000 \u000b2\b\u0012\u0004\u0012\u00020\u00000\u0001:\u0001\u000bB\u0011\b\u0002\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005R\u0010\u0010\u0002\u001a\u00020\u00038\u0006X\u0087\u0004¢\u0006\u0002\n\u0000j\u0002\b\u0006j\u0002\b\u0007j\u0002\b\bj\u0002\b\tj\u0002\b\n¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/EventQueuePriority;", "", "value", "", "<init>", "(Ljava/lang/String;II)V", "A", "B", TokenNames.f32C, "D", TokenNames.f33E, "Companion", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final class EventQueuePriority {
    private static final /* synthetic */ EnumEntries $ENTRIES;
    private static final /* synthetic */ EventQueuePriority[] $VALUES;

    /* renamed from: Companion, reason: from kotlin metadata */
    @InterfaceC88464a
    public static final Companion INSTANCE;

    @JvmField
    public final int value;

    /* renamed from: A */
    public static final EventQueuePriority f333481A = new EventQueuePriority("A", 0, 1);

    /* renamed from: B */
    public static final EventQueuePriority f333482B = new EventQueuePriority("B", 1, 2);

    /* renamed from: C */
    public static final EventQueuePriority f333483C = new EventQueuePriority(TokenNames.f32C, 2, 3);

    /* renamed from: D */
    public static final EventQueuePriority f333484D = new EventQueuePriority("D", 3, 4);

    /* renamed from: E */
    public static final EventQueuePriority f333485E = new EventQueuePriority(TokenNames.f33E, 4, 5);

    @Metadata(m64929d1 = {"\u0000\u0018\n\u0002\u0018\u0002\n\u0002\u0010\u0000\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\u0003\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0010\u0010\u0004\u001a\u0004\u0018\u00010\u00052\u0006\u0010\u0006\u001a\u00020\u0007¨\u0006\b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/EventQueuePriority$Companion;", "", "<init>", "()V", "findByValue", "Lcom/x/dmv2/thriftjava/EventQueuePriority;", "value", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class Companion {
        public /* synthetic */ Companion(DefaultConstructorMarker defaultConstructorMarker) {
            this();
        }

        @InterfaceC88465b
        public final EventQueuePriority findByValue(int value) {
            if (value == 1) {
                return EventQueuePriority.f333481A;
            }
            if (value == 2) {
                return EventQueuePriority.f333482B;
            }
            if (value == 3) {
                return EventQueuePriority.f333483C;
            }
            if (value == 4) {
                return EventQueuePriority.f333484D;
            }
            if (value != 5) {
                return null;
            }
            return EventQueuePriority.f333485E;
        }

        private Companion() {
        }
    }

    private static final /* synthetic */ EventQueuePriority[] $values() {
        return new EventQueuePriority[]{f333481A, f333482B, f333483C, f333484D, f333485E};
    }

    static {
        EventQueuePriority[] eventQueuePriorityArr$values = $values();
        $VALUES = eventQueuePriorityArr$values;
        $ENTRIES = EnumEntriesKt.m65223a(eventQueuePriorityArr$values);
        INSTANCE = new Companion(null);
    }

    private EventQueuePriority(String str, int i, int i2) {
        this.value = i2;
    }

    @InterfaceC88464a
    public static EnumEntries getEntries() {
        return $ENTRIES;
    }

    public static EventQueuePriority valueOf(String str) {
        return (EventQueuePriority) Enum.valueOf(EventQueuePriority.class, str);
    }

    public static EventQueuePriority[] values() {
        return (EventQueuePriority[]) $VALUES.clone();
    }
}